/* Development-only client for the ad-hoc XPC authority/descriptor probe. */

#include <errno.h>
#include <fcntl.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>
#include <xpc/xpc.h>

#ifndef BUILD
#define BUILD "unspecified"
#endif

static const char service_name[] = "dev.capsule.gate-b.license-free";
static const char expected_bytes[] = "exact-cross-process-xpc-bytes";

int main(int argc, char **argv) {
    bool malformed = argc == 2 && strcmp(argv[1], "--malformed") == 0;
    char path[] = "/tmp/capsule-gate-b-xpc.XXXXXX";
    int writable = mkstemp(path);
    if (writable < 0 ||
        write(writable, expected_bytes, sizeof(expected_bytes) - 1) !=
            (ssize_t)(sizeof(expected_bytes) - 1) ||
        close(writable) != 0) {
        unlink(path);
        return 1;
    }
    int content = open(path, O_RDONLY | O_NOFOLLOW | O_CLOEXEC);
    unlink(path);
    if (content < 0) {
        return 1;
    }

    xpc_connection_t connection = xpc_connection_create_mach_service(
        service_name, NULL, 0);
    xpc_connection_set_event_handler(connection, ^(xpc_object_t event) {
        (void)event;
    });
    xpc_connection_activate(connection);
    xpc_object_t request = xpc_dictionary_create_empty();
    xpc_dictionary_set_string(
        request, "operation", malformed ? "forged-operation" : "transfer-input");
    xpc_dictionary_set_string(request, "build", BUILD);
    xpc_dictionary_set_fd(request, "content", content);
    close(content);

    dispatch_semaphore_t reply_ready = dispatch_semaphore_create(0);
    __block xpc_object_t reply = NULL;
    xpc_connection_send_message_with_reply(
        connection, request, NULL, ^(xpc_object_t response) {
            reply = xpc_retain(response);
            dispatch_semaphore_signal(reply_ready);
        });
    xpc_release(request);
    if (dispatch_semaphore_wait(
            reply_ready, dispatch_time(DISPATCH_TIME_NOW, 5 * NSEC_PER_SEC)) != 0) {
        printf("xpc.timeout=true build=%s\n", BUILD);
        xpc_connection_cancel(connection);
        xpc_release(connection);
        return 3;
    }
    if (reply == NULL || xpc_get_type(reply) == XPC_TYPE_ERROR) {
        const char *description = reply == NULL ? "no reply" :
            xpc_dictionary_get_string(reply, XPC_ERROR_KEY_DESCRIPTION);
        printf("xpc.peer-denied=true build=%s reason=%s\n", BUILD,
               description == NULL ? "unknown" : description);
        if (reply != NULL) {
            xpc_release(reply);
        }
        xpc_connection_cancel(connection);
        xpc_release(connection);
        return 2;
    }

    int64_t status = xpc_dictionary_get_int64(reply, "status");
    bool identity = xpc_dictionary_get_bool(reply, "messageDerivedIdentityValid");
    bool read_only = xpc_dictionary_get_bool(reply, "fdReadOnly");
    bool bytes_exact = xpc_dictionary_get_bool(reply, "bytesExact");
    const char *reason = xpc_dictionary_get_string(reply, "reason");
    printf("xpc.status=%lld build=%s identity=%s fdReadOnly=%s bytesExact=%s reason=%s\n",
           (long long)status, BUILD, identity ? "true" : "false",
           read_only ? "true" : "false", bytes_exact ? "true" : "false",
           reason == NULL ? "none" : reason);
    xpc_release(reply);
    xpc_connection_cancel(connection);
    xpc_release(connection);
    if (malformed) {
        return status == 10 ? 0 : 5;
    }
    return status == 0 && identity && read_only && bytes_exact ? 0 : 6;
}
