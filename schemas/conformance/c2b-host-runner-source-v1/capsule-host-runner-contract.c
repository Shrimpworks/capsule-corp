/*
 * Capsule C2B host-runner source contract v1.
 *
 * PASSIVE STATIC EVIDENCE ONLY. This file is not a final runner artifact, is
 * not linked to libkrun, and must not be used to launch a VM or guest. It
 * freezes the C17 placement, preflight, descriptor, call-order, and authority
 * boundary that a later separately built and reviewed runner must preserve.
 */

#define CAPSULE_HOST_RUNNER_SOURCE_CONTRACT_VERSION 1
#define CAPSULE_RUNNERS_PER_ATTEMPT_ID 1
#define CAPSULE_EXECUTE_TIME_REPLACEMENT_VALUES 0
#define CAPSULE_HOST_CLOSE_FROM_INCLUSIVE 8
#define CAPSULE_ROOT_FD 4
#define CAPSULE_CONSOLE_PORT_COUNT 3

struct capsule_host_fd_contract {
    unsigned int fd;
    const char *role;
    const char *access_mode;
};

static const struct capsule_host_fd_contract capsule_host_fds[] = {
    {0, "runner-stdin-null", "O_RDONLY"},
    {1, "runner-stdout", "O_WRONLY"},
    {2, "runner-stderr", "O_WRONLY"},
    {3, "record-before-start-control", "O_RDONLY"},
    {4, "runtime-root", "O_RDONLY"},
    {5, "registered-source-port-input", "O_RDONLY"},
    {6, "approved-inline-input-port-input", "O_RDONLY"},
    {7, "completion-result-port-output", "O_WRONLY"},
};

struct capsule_console_port_contract {
    unsigned int port_id;
    const char *name;
    int input_fd;
    int output_fd;
    const char *guest_node;
};

static const struct capsule_console_port_contract capsule_console_ports[] = {
    {0, "capsule.source", 5, -1, "/dev/hvc0"},
    {1, "capsule.input", 6, -1, "/dev/vport0p1"},
    {2, "capsule.completion", -1, 7, "/dev/vport0p2"},
};

static const char *const capsule_virtio_devices[] = {
    "balloon",
    "rng",
    "console-multiport-with-three-fixed-ports",
    "block-root-vda-read-only",
};

static const char *const capsule_forbidden_calls[] = {
    "krun_set_kernel",
    "krun_set_firmware",
    "krun_add_disk",
    "krun_set_root",
    "krun_add_virtiofs",
    "krun_add_virtiofs2",
    "krun_add_virtiofs3",
    "krun_add_virtiofs4",
    "krun_add_vsock",
    "krun_add_vsock_port",
    "krun_add_vsock_port2",
    "krun_add_net_unixstream",
    "krun_set_gpu_options",
    "krun_set_snd_device",
};

/* Each line below is an ordered, fail-closed source obligation. */
#define CAPSULE_PREFLIGHT(order, operation) \
    static const char *const capsule_step_##order = "preflight:" #operation
#define CAPSULE_LIBKRUN_CALL(order, operation) \
    static const char *const capsule_step_##order = "libkrun:" #operation
#define CAPSULE_HANDSHAKE(order, operation) \
    static const char *const capsule_step_##order = "handshake:" #operation

CAPSULE_PREFLIGHT(01, validate_one_supervisor_owned_process_per_attempt_id);
CAPSULE_PREFLIGHT(02, validate_exact_argv0_empty_environment_fds_0_through_7_close_from_8_no_execute_time_replacements);
CAPSULE_PREFLIGHT(03, validate_fd4_unlinked_regular_mode_0400_ordonly_fixed_length_and_digest);
CAPSULE_LIBKRUN_CALL(04, krun_create_ctx);
CAPSULE_LIBKRUN_CALL(05, krun_set_vm_config_vcpus_1_ram_mib_256);
CAPSULE_LIBKRUN_CALL(06, krun_disable_implicit_console);
CAPSULE_LIBKRUN_CALL(07, krun_disable_implicit_init);
CAPSULE_LIBKRUN_CALL(08, krun_disable_implicit_vsock);
CAPSULE_LIBKRUN_CALL(09, krun_add_read_only_raw_root_fd_fd4_vda);
CAPSULE_LIBKRUN_CALL(10, krun_set_root_disk_remount_dev_vda_ext4_ro_nosuid_nodev);
CAPSULE_LIBKRUN_CALL(11, krun_add_virtio_console_multiport_require_console_id_0);
CAPSULE_LIBKRUN_CALL(12, krun_add_console_port_inout_capsule_source_fd5_minus1_require_port_id_0);
CAPSULE_LIBKRUN_CALL(13, krun_add_console_port_inout_capsule_input_fd6_minus1_require_port_id_1);
CAPSULE_LIBKRUN_CALL(14, krun_add_console_port_inout_capsule_completion_minus1_fd7_require_port_id_2);
CAPSULE_LIBKRUN_CALL(15, krun_set_kernel_console_hvc0);
CAPSULE_LIBKRUN_CALL(16, krun_set_workdir_root);
CAPSULE_LIBKRUN_CALL(17, krun_set_exec_fixed_init_single_argv_three_entry_environment);
CAPSULE_HANDSHAKE(18, emit_ready_then_require_exact_one_byte_G_followed_by_eof);
CAPSULE_LIBKRUN_CALL(19, krun_start_enter);
