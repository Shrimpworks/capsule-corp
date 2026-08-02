# P0-1 research handoff

- Question: can pinned libkrun consume a finalized unlinked runtime root through direct inherited
  read-only descriptor custody without same-user substitution?
- Defensive scope: owned repository fixtures, local processes, cached pinned OCI images, and one
  owned local libkrun/HVF guest only.
- Decision: **PATCH-REQUIRED**; use a narrow raw-only FD-native libkrun API and rerun the final
  signed installed App Sandbox corpus. Do not add a privileged helper.
- Confidence: high for the observed unsandboxed two-open identity, descriptor semantics, local
  alias/mapping negatives, positional-I/O source path, and guest digest; low for installed App
  Sandbox/end-to-end same-UID custody because no valid code-signing identity was available.
- Primary evidence: `RESULTS.md` and `evidence/2026-08-02/`.
- Exact residual test: final Developer ID/notarized Supervisor protected-container construction,
  direct runner inheritance, App Sandbox `/dev/fd`/FD-native attachment, task-port/grant denial,
  crash/recovery, and guest digest on exact final bytes.
- Prototype disposition: retain as development-only until the patched installed corpus is
  reconciled; product packages must not import it.
