# Gate B macOS Authority-Separation Spike

Status: development-only disposable experiment. Nothing here is a production security boundary or
receipt evidence.

Owner: Capsule architecture spike. Remove or replace this experiment after independently
distribution-signed daemon, Broker, and Supervisor targets have passed the Gate B matrix on the
minimum supported macOS release and the resulting contracts are frozen.

## Bounded question

Can macOS make the daemon, Trusted Host Broker, and Execution Supervisor distinct authorities for
IPC, keys, and storage, rather than relying only on same-user process packaging?

## Hypothesis

Apple-issued component signatures plus XPC peer requirements, data-protection Keychain access
groups, Secure Enclave access controls, and protected app/app-group containers can enforce most of
the separation. Exact trust epochs and migration remain protocol responsibilities. Development or
ad-hoc signing cannot substantiate the production boundary.

## Reproduction

Run `./run.sh --with-debugger` on macOS. The script creates only derived files under `build/`. The
key probe creates ephemeral Keychain keys/items and deletes them before exit. It deliberately
suppresses interactive authentication; it must not satisfy the approval-key user-presence policy.

The retained source covers:

- correct and wrong signing identifiers;
- an unsigned binary;
- an ad-hoc impostor with the expected identifier;
- two builds with the same identifier;
- an exact copied binary;
- rejection of ad-hoc code by an Apple-chain requirement;
- rejection of the development-only `get-task-allow` entitlement;
- point-in-time dynamic-valid and sticky debugger-attached status;
- an explicit unentitled data-protection Keychain access group;
- Secure Enclave P-256 creation, background signing, a persistent-key export branch when
  entitlement-backed persistence succeeds, and noninteractive user-presence denial.

It does not claim to be an end-to-end XPC test. That requires three Apple Development or
distribution identities, validated entitlements/provisioning profiles, installed launchd/XPC
services, protected containers, and interactive user presence.

See [RESULTS.md](RESULTS.md) for the evidence/limitation split, Gate B decision, proposed document
changes, and next test.
