#include <errno.h>
#include <mach/mach.h>
#include <mach/mach_error.h>
#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <sys/mman.h>
#include <unistd.h>

static int set_limit(int megabytes, const char *label) {
    int old_limit = 0;
    kern_return_t result = task_set_phys_footprint_limit(
        mach_task_self(), megabytes, &old_limit);
    fprintf(stderr, "%s result=%d (%s) old_limit_mb=%d\n",
        label, result, mach_error_string(result), old_limit);
    return result == KERN_SUCCESS ? 0 : 1;
}

static int touch_megabytes(size_t megabytes) {
    size_t length = megabytes * 1024U * 1024U;
    unsigned char *mapping = mmap(NULL, length, PROT_READ | PROT_WRITE,
        MAP_PRIVATE | MAP_ANON, -1, 0);
    if (mapping == MAP_FAILED) {
        fprintf(stderr, "mmap bytes=%zu errno=%d (%s)\n",
            length, errno, strerror(errno));
        return 1;
    }
    for (size_t offset = 0; offset < length; offset += (size_t)getpagesize()) {
        mapping[offset] = 1;
    }
    fprintf(stderr, "touched_mb=%zu\n", megabytes);
    return munmap(mapping, length) == 0 ? 0 : 1;
}

int main(int argc, char **argv) {
    if (argc != 3 || (strcmp(argv[1], "hold") != 0 && strcmp(argv[1], "raise") != 0)) {
        fprintf(stderr, "usage: %s hold|raise TOUCH_MEGABYTES\n", argv[0]);
        return 64;
    }
    char *end = NULL;
    unsigned long touch = strtoul(argv[2], &end, 10);
    if (end == argv[2] || *end != '\0' || touch > 512UL) {
        return 64;
    }
    if (set_limit(128, "set-128") != 0) {
        return 65;
    }
    if (strcmp(argv[1], "raise") == 0 && set_limit(512, "raise-512") != 0) {
        return 66;
    }
    return touch_megabytes((size_t)touch);
}
