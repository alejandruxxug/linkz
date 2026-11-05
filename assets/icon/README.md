# The Linkz icon

The approved icon (R006, option E2): `linkz-1024.png` is the master and `linkz-16.png` the hand-drawn 16 px sprite. Both are copies of `docs/Art/Reviews/R006/E2-cool-brush-in-teeth-one-eye-{1024,16}.png`; the generator is `docs/Art/Figma-gen/r006c.py` (run it, then copy the two files here again).

`scripts/crossbuild.sh` passes them to `lizardscript build --icon --icon-small`; the editor and the splash use the master as the window's icon (`app_icon()` in `lib/app/editor.liz`).
