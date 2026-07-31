# Retained observed results

Date: 2026-07-31 (America/Toronto)

Repository baseline: `9bfd2acedbccfbe851f797edc06eb447733188e3`

Upstream pins:

- `apple/container` 1.0.0: `ee848e3ebfd7c73b04dd419683be54fb450b8779`
- `apple/containerization` 0.33.3: `a2a1add6c7e1a1665e5397edc49d925c49090b3a`

## Local smoke output

```text
{"language":"swift","malformedCodeRequirementStatus":22,"sameTeamRequirementStatus":0,"selfDynamicCodeValidityStatus":0,"validCodeRequirementStatus":0}
binary_bytes=73000 path=.../.build/swift-platform-probe
ok capsule.local/capsule/experiments/gate-e-supervisor-topology/go-platform-probe 0.275s
{"language":"go-cgo","validCodeRequirementStatus":0,"malformedCodeRequirementStatus":22,"sameTeamRequirementStatus":0,"selfDynamicCodeValidityStatus":0}
binary_bytes=1911378 path=.../.build/go-platform-probe
```

Peak RSS from `/usr/bin/time -l`:

```text
Swift:  7,880,704 bytes
Go:    10,862,592 bytes
```

## Apple backend authorization probe

Observed launchd shape after temporary start:

```text
domain: gui/501
label: com.apple.container.apiserver
type: LaunchAgent
program: /usr/local/bin/container-apiserver
```

Ad-hoc Go XPC client outside the Codex sandbox:

```json
{"language":"go-cgo","validCodeRequirementStatus":0,"malformedCodeRequirementStatus":22,"sameTeamRequirementStatus":0,"selfDynamicCodeValidityStatus":0,"appleApiPingStatus":0}
```

Negative control from inside the Codex sandbox: `appleApiPingStatus=2` (XPC error).

Final service status:

```text
apiserver is not running and not registered with launchd
```

## Source and installed-binary measurements

```text
PASS exact Apple sources retain the tested privilege and client-authentication shape

59,906,384 /usr/local/bin/container
55,197,520 /usr/local/bin/container-apiserver
54,879,744 /usr/local/libexec/container/plugins/container-runtime-linux/bin/container-runtime-linux
54,474,720 /usr/local/libexec/container/plugins/container-core-images/bin/container-core-images
52,916,288 /usr/local/libexec/container/plugins/container-network-vmnet/bin/container-network-vmnet

20,971 lines apple/containerization Sources/Containerization
50,920 lines apple/containerization Sources
40,893 lines apple/container Sources
36 routes in apple/container 1.0.0 XPCRoute
```

These results are spike evidence only. They are not backend validation, platform attestation, or a
production security claim.

