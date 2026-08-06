---
title: "Capsule’s First Hostile Guest: From Sandbox Design to a Real Isolation Test"
description: "How Capsule booted a disposable micro-VM, ran a governed Deno workload, and proved a fixed set of dangerous capabilities stayed unavailable."
date: 2026-08-06
status: draft
---

# Capsule’s First Hostile Guest: From Sandbox Design to a Real Isolation Test

For a long time, Capsule’s security architecture existed in the place where a
lot of ambitious systems begin: diagrams, schemas, threat models, reproducible
builds, and small experiments proving individual parts could work.

That work matters. But eventually a sandbox has to do the uncomfortable thing:
run something that actively tries to break its rules.

We have now reached that point.

In our latest experiment, Capsule launched one real, disposable Linux guest on
an Apple-silicon Mac using libkrun and Apple’s Hypervisor framework. Inside that
guest, a fixed Deno-based workload first proved that useful code could run. It
then switched to a deliberately constrained identity and attempted a collection
of operations that Capsule intends to forbid.

All 30 expected test markers matched. The guest could not write to its root
filesystem, mount a filesystem, regain root, reach host-only paths, obtain a
usable vsock transport, open a raw block device for writing, find an active
non-loopback network path, access virtiofs, or inherit an ambient environment.
The process exited normally, the Supervisor reaped it, the disposable root was
fully unlinked, and a host-side canary remained unchanged.

This is not a production-security claim, and it is not an alpha release. It is
something narrower and, for the project, extremely important:

**the first end-to-end evidence that Capsule’s intended authority boundary can
run real code and withstand a fixed hostile test inside a real guest.**

## Why put a sandbox inside a guest?

Deno and `deno_core` offer excellent tools for constructing a restricted
JavaScript runtime. We use a custom snapshot, a deliberately small operation
surface, no general module loader, no inspector, disabled string code
generation, and other runtime restrictions.

Those controls are useful, but an in-process JavaScript sandbox is not the host
security boundary.

If the runtime, a native dependency, or some future integration contains a
memory-safety bug, the process around the JavaScript engine may be compromised.
A serious sandbox therefore needs another boundary outside the runtime—one that
does not depend on every line of the runtime being perfect.

Capsule’s design places untrusted execution inside a small disposable guest. The
trusted application remains on the host and keeps the authority it needs to
manage the product. The guest receives only the exact artifacts and resources
approved for one attempt.

The shape is roughly:

```text
Trusted Capsule application
        |
        v
Approval + exact execution plan
        |
        v
Execution Supervisor (sole guest owner)
        |
        v
Exact signed runner -> libkrun/HVF
        |
        v
Disposable Linux guest
        |
        v
Fixed launcher -> governed deno_core workload
```

This separation lets the app remain a capable, full-featured product without
granting the same authority to the code it executes. The design goal is not
“make the whole app untrusted.” It is “keep untrusted work away from the app’s
authority.”

## The first breakthrough: useful code actually ran

Before testing hostile behavior, we needed to prove the basic execution path was
real.

That happened in version 19 of the experiment.

The Execution Supervisor verified the exact profile and every input artifact,
created a private attempt root, opened it, and removed its pathname before
launch. It started the exact runner with an empty environment and a sealed set
of descriptors, verified the live runner’s identity, persisted an attempt
record, and only then issued a one-use start authorization.

Inside the guest, a small trusted PID 1 prepared a read-only environment and
launched one fixed `deno_core` program. The program produced the expected known
answer and evidence that the runtime’s host-side syscall seal was active.

The complete console proof matched exactly. The runner exited normally, the
Supervisor observed `waitpid`, no forced kill was necessary, and teardown
completed in five milliseconds.

The whole authorization-to-result interval was 221 milliseconds.

That v19 result answered a foundational question: could the Supervisor,
libkrun/HVF, a read-only Linux root, a trusted launcher, and Capsule’s governed
Deno runtime all work together in one real attempt?

Yes.

But the workload was deliberately benign and fixed. A successful “hello world”
inside a VM does not tell you whether the isolation policy holds when the guest
pushes against it.

## What “hostile guest” means here

“Hostile” can sound dramatic, so it is worth being precise.

This was not an uncontrolled exploit kit, a network scanner, or arbitrary code
downloaded from somewhere else. It was a small, fixed, auditable test program
inside an owned disposable guest. It attempted only named local operations that
the Capsule profile is supposed to deny or make unavailable.

The test had no credential authority. It had no network authority. Its source
contained no `connect`, `send`, `sendto`, or `sendmsg` call. It did not target
another machine or identity.

The point was to ask a bounded question:

> If code inside this exact guest behaves adversarially, do the specific
> boundaries we configured produce the denials and absences we expect?

The hostile phase checked that the guest:

- ran as a non-root PID 1;
- had `no_new_privs` enabled and zero effective Linux capabilities;
- received only the expected file descriptors;
- could not modify the read-only root filesystem;
- could not see selected host-only paths;
- could not mount another filesystem;
- could not regain uid 0;
- could not obtain a usable local-CID vsock transport;
- could not open its block device for writing;
- had no virtio network device, active non-loopback interface, or non-loopback
  IPv4/IPv6 route;
- had no virtiofs mount; and
- received an empty ambient environment.

Every unexpected success was defined as a failed attempt. A normal process exit
alone was not enough to pass; the Supervisor required the complete exact
transcript plus its own independent lifecycle and teardown evidence.

## The failures were part of the result

The final passing version was v27. We did not jump directly from v19 to v27,
and the versions in between are one of the most valuable parts of the work.

Each attempt was bound to exact bytes and a one-use authorization. When an
attempt stopped, we retained the evidence, diagnosed the narrowest supported
cause, produced a new immutable candidate, and required another authorization.
We did not quietly weaken the check or reinterpret partial output as success.

Here is the shorter version of that journey:

| Version | What happened | What we learned |
| --- | --- | --- |
| v20-v22 | The runner refused before guest authorization. | The fail-closed path worked, but early error evidence needed to converge reliably. |
| v23 | Extra hashing evidence showed that the real root was correct. | The expected digest embedded in the runner had been malformed. The guest still never launched. |
| v24 | The corrected guest passed identity, descriptor, root, host-path, mount, and privilege checks. | The next ambiguity was in the vsock test. |
| v25 | We stopped before launch. | Socket creation alone does not prove usable transport, so the proposed test measured the wrong property. |
| v26 | The active vsock and raw-block tests passed, then network inventory failed. | The guest had an expected down, unbacked `dummy0`; rejecting every interface name other than `lo` was too crude. |
| v27 | The network check was rewritten around device backing, flags, and routes. | The complete 30-marker corpus passed. |

Two moments are especially important.

First, v23 proved the value of independent evidence. The Supervisor hashed the
root by pathname before unlinking it, hashed the same open descriptor after
unlinking it, and compared those values with the digest computed by the runner.
All three agreed. The failure was not mysterious hypervisor corruption; it was
an incorrectly embedded expected byte array. The system refused to launch
because its bytes did not match its own declaration.

Second, v25 was marked **NO_GO without running it**. Review showed that its
vsock predicate confused “a socket object can be created” with “a transport is
usable.” Launching another guest would have produced more output but not better
knowledge. We corrected the question first.

This is the kind of iteration security work needs: preserve adverse evidence,
improve observability, and test the intended property rather than chasing a
green result.

## The v27 result

The final v27 candidate kept the previously validated controls and corrected
the passive network check.

It allowed only loopback or loopback plus the known `dummy0` shape. If `dummy0`
existed, it had to be down and have no backing device. The probe rejected a
virtio-net modalias, any IPv4 default route, and every non-loopback IPv4 or IPv6
route. The separate active vsock check still required local-CID transport to be
unavailable, without connecting or sending.

One exact, user-authorized Supervisor invocation produced:

- 30 console lines and 30 exact expected markers;
- the final marker `C2BHOSTILE20:complete`;
- exact completion and console proofs;
- runner exit status zero and authoritative `waitpid` observation;
- no forced kill;
- 208 milliseconds from authorization to result;
- five milliseconds of teardown;
- an unlinked guest root with final link count zero;
- matching root hashes before and after unlink; and
- an unchanged host canary.

The preflight record also retained:

```text
networkAuthorized: false
credentialsAuthorized: false
externalConnectionAttempted: false
networkBytesSent: false
maximumGuests: 1
```

That distinction matters. The guest’s own “complete” marker was only one part
of the proof. Capsule also retained host-side evidence about what was
authorized, which process ran, how it exited, whether it was reaped, what
happened to the root, and whether the canary changed.

In other words, the guest did not get to grade its own exam.

## What we can honestly claim now

For this exact materialized profile, on this owner-controlled Apple-silicon Mac,
in one disposable local reproduction, we have evidence that:

- the Supervisor can bind and authorize one exact guest attempt;
- the selected libkrun/HVF path can boot it;
- immutable inputs can reach a governed Deno workload without a live host path;
- the expected workload can complete;
- the fixed hostile corpus observes the expected denials and absences;
- the runner can be reaped normally; and
- the attempt root and host canary can finish in the expected state.

This eliminates a major feasibility uncertainty. We are no longer asking
whether the core pieces can compose at all. They did.

## What this does not prove

The easiest mistake at a checkpoint like this is to turn a narrow result into a
broad slogan. “The hostile test passed” is accurate. “Capsule is secure” is not.

This experiment does not yet prove:

- safe execution of arbitrary user-supplied JavaScript;
- the final typed source, input, and completion transport;
- the installed app’s complete authentication, approval, signing, recovery,
  and one-use grant path;
- production packaging, Developer ID signing, notarization, Gatekeeper, or
  clean-host installation;
- behavior under memory or CPU pressure, sleep/wake, upgrades, crashes,
  response loss, repetition, concurrency, or long-running soak tests;
- containment against every possible guest-kernel, hypervisor, or
  microarchitectural attack; or
- alpha or production admission.

Passive device and route inventory is also not mathematical proof that every
possible networking path is absent. It is evidence for the exact controls we
tested.

The correct project status remains **IN_PROGRESS — TRENDING_GOOD**.

## Why this still feels like a turning point

There is a meaningful difference between an architecture that should work and
one that has survived its first real confrontation with reality.

We now have a chain of evidence from deterministic materialization, through
one-use authorization, into a real guest, through a useful Deno execution, into
adversarial local checks, and back out through process ownership and teardown.
The failures along the way strengthened that chain because they showed that
incorrect identities, malformed digests, ambiguous checks, and policy mismatches
did not get waved through.

For me, the biggest milestone is not simply that version 27 turned green. It is
that Capsule is beginning to behave like the system we designed:

- authority stays on the trusted side;
- execution receives only a narrow, immutable plan;
- starting a guest is an explicit, one-use act;
- partial evidence does not become success;
- the guest cannot declare itself trustworthy; and
- cleanup is part of correctness, not an afterthought.

## What comes next

The next phase is less about proving that a guest can exist and more about
turning this experimental path into a governed product path.

The major remaining work includes:

1. closing the final typed source/input/completion transport;
2. wiring authenticated submission and human approval to the installed
   Supervisor without moving authority into the daemon;
3. admitting arbitrary approved source only after the relevant hostile source
   and transport corpus passes;
4. validating cancellation, crashes, response loss, pressure, recovery,
   sleep/wake, and upgrades;
5. testing installed, signed, and notarized artifacts on clean supported hosts;
   and
6. completing independent security and release admission.

That is still a substantial list. But it is a different list than we had before
this checkpoint. The question is no longer “can Capsule run a governed workload
inside a tightly controlled guest?” The question is how carefully we can carry
this proven experimental shape into the real product.

And that is a very good place for the project to be.

## Technical receipts

For readers who want the exact identities:

### Benign execution, v19

- Attempt:
  `capsule-c2b-v19-immutable-fixture-benign-owned-guest-20260806-01`
- Composed-profile digest:
  `ac2721719a1e4f15c664e0b7c21d99602b6fc7d5a9c55c8b17d08970098f48fa`
- Result: **PASSED**
- Authorization to result: 221 ms
- Teardown: 5 ms

### Hostile denial execution, v27

- Attempt:
  `capsule-c2b-v27-hostile-network-corrected-owned-guest-20260806-01`
- Composed-profile digest:
  `52f38c8f964a59dbf7e7ed98576ee95aae0470cba2462749551e7b335ca6073e`
- Guest-root SHA-256:
  `002524fb0cf1b03df110bbb8c243751cf259f50e19dd85bf84f52ce30d80119d`
- Signed runner SHA-256:
  `49127899025f1216cfdafd54079557967d8a5c677fddbbe23c5d6bef0230f86b`
- Kernel-console SHA-256:
  `b0a593750065500c99a193bff62de43992324d47508dc4daadf0a827e7181f74`
- Result: **PASSED**
- Authorization to result: 208 ms
- Teardown: 5 ms

The full engineering checkpoint, including the v10-v27 progression and claim
boundaries, is available in
[`FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md`](https://github.com/Shrimpworks/capsule-corp/blob/526b300f59674afc727ce80dcc106603a3503a07/docs/FIRST_OWNED_GUEST_EXECUTION_CHECKPOINT.md).
