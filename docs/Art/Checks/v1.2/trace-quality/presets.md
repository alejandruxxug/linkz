# Trace quality presets: measured

`lib/trace/trace_quality.liz`. Measured with `scripts/trace_score.sh --quality NAME` (or `--set key=value`), the yardstick of `lib/trace/trace_measure.liz`: mean color error (0 to 255), ink IoU, code size.

## Presets

| preset | options over the defaults |
|---|---|
| light | colors 14, detail 1.4, line_min 10, snap 2.0, gradient_error 7 |
| balanced | none (the defaults, byte-identical) |
| detailed | colors 32, detail 0.4, line_min 3, snap 0.8, gradient_error 6 |
| max | detailed + max_side 1600 |

## 4 references (wikipe-tan-casual, clean-kamina, messy-pleroma-fox, wikipe-tan-conducting)

| preset | mean err | mean IoU | code | note |
|---|---:|---:|---:|---|
| light | 7.06 (+6.5%) | 0.328 | 0.39 MB (-42%) | final values, with gradient_error 7 |
| balanced | 6.63 | 0.334 | 0.67 MB | |
| detailed | 6.31 (-4.8%) | 0.339 | 1.41 MB | final values, with gradient_error 6 |
| max | 5.87 (-11.5%) | 0.358 | 2.08 MB | only wikipe-tan-casual and messy-pleroma-fox are bigger than 1024 px |

Light candidates tried (err, code): colors 10/detail 1.3/line_min 8/snap 2.0: 8.35, 0.42 MB; colors 8/1.6/10/2.5: 9.35, 0.34 MB; colors 12/1.2/8/1.5: 7.21, 0.46 MB; 16/1.5/10/2.0: 6.85, 0.45 MB; 14/1.3/8/1.5: 6.94, 0.47 MB; 12/1.5/10/1.5: 7.28, 0.40 MB; 14/1.4/10/2.0 (chosen, then gradient_error 7): 7.00, 0.45 MB. Colors drive the error; the code shrinks mostly with line_min, snap and detail.

## gradient_error, on 4 messy references (hyou-chan, pleroma-fox, pepper-portrait, coriander-concept)

| preset | gradient_error | mean err | mean IoU | code |
|---|---:|---:|---:|---:|
| balanced | 4 (today) | 5.67 | 0.208 | 1.56 MB |
| balanced | 6 | 5.74 | 0.212 | 1.47 MB |
| balanced | 8 | 5.77 | 0.209 | 1.46 MB |
| detailed | 4 | 5.03 | 0.231 | 4.04 MB |
| detailed | 6 (chosen) | 5.04 | 0.231 | 3.99 MB |
| detailed | 8 | 5.04 | 0.228 | 3.92 MB |
| light | 4 | 6.00 | 0.201 | 0.95 MB |
| light | 7 (chosen) | 6.06 | 0.203 | 0.89 MB |
| light | 10 | 6.06 | 0.202 | 0.89 MB |

Stronger gradients cost nothing in the preset that has many colors (detailed) and a little (+1%) in light; the saving is in code (2 to 6%, and 38% on a flat picture like messy-pleroma-fox). Balanced keeps 4 because it must stay byte-identical.

## Auto's score

`score = err + 0.5 * MB of code` (`quality_mb_cost`), lowest wins; ties go to the first of light, balanced, detailed. A megabyte of code has to buy half a color level of error (about 8 to 10% of a typical error).

* Clean picture (wikipe-tan-conducting): detailed is -0.67 err for +0.17 MB: 0.67 > 0.085, detailed wins.
* The four references above: detailed beats balanced by 1.0 to 3.9 err per MB, well over 0.5, so it wins on all of them; light never beats balanced (3.4 to 12 err per MB saved).
* Big painting: 3 MB more for 0.2 less error costs 1.5, so it loses; 2 MB less for 0.4 more error (light) wins.

Auto's rows go in the stats under `picked.rows`, with `picked.less_error` (percent of balanced's error) and `picked.more_mb`.
