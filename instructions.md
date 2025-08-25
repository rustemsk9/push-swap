# Push_swap: Practical Strategy and Engineering Notes

Goal
- Produce short instruction sequences to sort stack A (ascending) using the limited ops: pa, pb, sa, sb, ss, ra, rb, rr, rra, rrb, rrr.
- Balance total operations and implementation complexity.

Constraints recap
- You can only:
  - push: pa, pb
  - swap top: sa, sb, ss
  - rotate up: ra, rb, rr
  - rotate down: rra, rrb, rrr
- No random access; everything is top-based and rotation-based.

Recommended high-level strategy by input size
- N ≤ 3: hardcoded minimal solutions (two/three).
- N ≤ 5: push smallest(s) to B, sort three on A, pa/pa inserting correctly, use ss when both benefit.
- 6 ≤ N ≤ ~120: Longest Increasing Subsequence (LIS) or chunking to decide what to push; Greedy insertion (“RotPlan”) to move elements with minimal cost; keep B strictly descending.
- N > ~120 (e.g., 500):
  - Easiest: Radix sort (LSD) on compressed ranks (0..n-1) using only pb/pa/ra; predictable and fast in practice.
  - Better counts: Larger-chunk LIS + RotPlan.

RotPlan: the insertion engine (core greedy step)
- Purpose: compute-and-execute the cheapest way to bring A[i] to top and place it into B at its correct spot (descending), aggregating rotations with rr/rrr.
- Steps:
  1) For candidate x = A[idxA], find insertion index idxB in B to preserve strict descending order.
  2) Compute rotation distances:
     - raUp = idxA, rraDown = lenA - idxA
     - rbUp = idxB, rrbDown = lenB - idxB
  3) Evaluate four patterns:
     - both up: cost = max(raUp, rbUp)
     - both down: cost = max(rraDown, rrbDown)
     - mixed 1: cost = raUp + rrbDown
     - mixed 2: cost = rraDown + rbUp
     - totalCost = patternCost + 1 (for pb)
  4) Pick minimal plan, execute:
     - While both need up: rr
     - While both need down: rrr
     - Finish remaining single rotations (ra/rb/rra/rrb)
     - pb (or pa when inserting into A)

- Target idxB in B (strictly descending):
  - If B empty: 0
  - Otherwise, scan circular pairs (prev,curr):
    - If prev >= x >= curr, insert at curr
  - If no slot found: insert after max(B) (x becomes new top or right below max).

Why RotPlan works well
- It compresses rotations (rr/rrr), avoiding naive separate moves.
- It turns insertion into a tractable local optimization that aligns with available ops.
- Combined with a good global selection (what to push, when), it yields strong totals.

Global selection strategy
- LIS approach (typical best for 100):
  - Rank-compress values (map to 0..n-1).
  - Compute an LIS (ascending) on A; mark LIS elements to keep in A.
  - pb all non-LIS to B while maintaining B descending (use rb when needed).
  - Then, while B not empty, insert best candidate using RotPlan into A (ascending).
  - Final cleanup: rotate A so min(A) is at top if required by your checker.

- Chunking approach (simpler to implement):
  - Partition ranks into K chunks (e.g., 5–8 for 100; ~11 for 500).
  - Process chunks in order: push elements of current chunk to B, keep B descending.
  - Insert back into A using RotPlan.

- Radix (LSD) baseline (great for 500+ due to simplicity):
  - On ranks, for each bit from LSB to MSB:
    - For i in 0..lenA-1: if bit = 0: pb; else: ra
    - Then pa until B empty
  - Very simple, predictable, passes performance requirements with modest counts.

Small-stack routines: keep them minimal and B-aware
- two(A,B):
  - If A[0] > A[1]: use ss if B[0] < B[1], else sa.
- three(A,B):
  - Cover 6 permutations with minimal {sa, ra, rra}. Opportunistically use ss if B benefits at the same time.
- four_Five(A,B):
  - Push the smallest (for 4) or two smallest (for 5) to B.
  - Sort three on A with the above routine.
  - pa / pa, using ra/rra/sb/rr/rrr opportunistically to place correctly.
  - Avoid extra pushes once you’re in the small path.

Engineering flow (end-to-end)
1) Parse, validate, deduplicate; compress ranks to 0..n-1.
2) If already sorted: print nothing and exit.
3) If n ≤ 5: run small-stack path.
4) Else if n ≤ ~120:
   - Compute LIS; pb all non-LIS to B (maintain B descending).
   - While B not empty: for each candidate in B (or examine a window), compute RotPlan into A; execute best; repeat.
   - Optional: when both tops suggest benefit, use ss; exploit rr/rrr aggregation.
5) Else:
   - Implement radix (LSD) for a robust baseline, or a larger-chunk LIS approach for better counts.

Expected results (ballpark, implementation-dependent)
- 100 elements:
  - LIS + RotPlan: ~700–1000 ops with good heuristics
  - Chunking: ~1000–1500
  - Radix: ~900–1300
- 500 elements:
  - LIS + RotPlan: ~5000–7500
  - Chunking: ~7000–11000
  - Radix: ~6500–9000

Implementation checklist
- Rank compression and safe parsing/errors.
- Core ops print their mnemonic (pb, pa, ra, etc.) exactly once per call.
- B-target index function (descending).
- RotPlan:
  - Rotation distances
  - Cost evaluation for 4 patterns (+1 push)
  - rr/rrr aggregation, then residual rotations, then push
- MaybeSS heuristic:
  - If A[0] > A[1] and B[0] < B[1], prefer ss over separate sa/sb when it reduces next plan cost.
- LIS or chunking pipeline (choose one first; radix for large as needed).
- Small-stack routines (two/three/four_Five) that never undo global work.

Testing guidance (macOS examples)
- Quick check:
  - ARG="4 67 3 87 23"; ./push_swap $ARG | ./checker $ARG
- Save ops and inspect:
  - ARG="2 1 3 6 5 8"; ./push_swap $ARG | tee ops.txt | ./checker $ARG
- Fuzz:
  - python - << 'PY'
    import random, subprocess
    for n in [5,10,50,100]:
        for _ in range(50):
            arr = random.sample(range(10000), n)
            args = ' '.join(map(str, arr))
            ps = subprocess.run(['./push_swap', *map(str,arr)], capture_output=True, text=True).stdout
            ck = subprocess.run(['./checker', *map(str,arr)], input=ps, capture_output=True, text=True).stdout.strip()
            if ck != 'OK': print('KO for', n, arr); raise SystemExit
    print('OK on samples')
  - PY

Convert this document to PDF
- VS Code: install “Markdown PDF”, open this file, Cmd+Shift+P → “Markdown PDF: Export (pdf)”
- Pandoc (Homebrew):
  - brew install pandoc
  - pandoc docs/push_swap_approach.md -o docs/push_swap_approach.pdf

References and learning
- Longest Increasing Subsequence (patience sorting): https://cp-algorithms.com/sequences/longest_increasing_subsequence.html
- Push_swap strategy writeups and tutorials:
  - https://medium.com/@ayogun/push-swap-c1f5d2d41e97
  - https://github.com/aitorres/push_swap_tutorial
- Radix sort (LSD) concept: https://en.wikipedia.org/wiki/Radix_sort

Summary
- Use RotPlan as your core local optimizer (the “insertion engine”) to minimize rotations and exploit rr/rrr.
- Pair it with LIS or chunking to decide what to push and when.
- For very large inputs, radix is the simplest reliable baseline.
- Keep small-stack routines minimal and avoid undoing global progress.