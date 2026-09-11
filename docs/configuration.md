# Configuration

`tele` keeps one YAML file, created on first run:

```text
~/.config/tele/config.yml
```

Pass `--config <path>` to use a different one.

## The settings overlay

Press `,` for the settings overlay: every setting `tele` has, grouped and ordered
the way the config file is, so a row you see there is found in the file at the
same path. It shows what each one is worth now, marks the ones nobody has chosen
(`[default]`), and says when a change takes hold - at once, `[next]` time `tele`
does that thing, or after a `[restart]`. Keybindings are listed there too, with
what your config changed and what it changed from.

Changes made there are written straight into `config.yml`, keeping your
comments, blank lines and anything in the file `tele` does not know about. Editing
the file in an editor and editing it in the overlay are the same act: both end
with the file being read again, so the two can never disagree. `enter` changes
the setting under the cursor, `r` puts it back to its default, and `esc` closes.

Four settings are shown but not changed there - which account this is, and where
its state lives. Changing those means moving files and reopening a session, so
they are a deliberate edit to the file rather than a keystroke.

## config.yml

```yaml
# state_dir: ~/.local/state/tele # session, local database and instance lock

ui:
  history_limit: 50 # messages fetched per chat on open
  # theme: # omit for the built-in tele-dark / tele-light
  #   dark: my-dark # ~/.config/tele/themes/my-dark.yml
  #   light: my-light

  notifications:
    desktop: true # hand it to the OS notification service
    toast: true # draw it in a corner of tele's own window
    preview: true # set false to send the sender's name and no message text

  toasts:
    error_zone: bottom-right # bottom-right | top-right
    notify_zone: top-right # bottom-right | top-right
    max_visible: 3 # per corner; the rest are counted, not drawn

photos:
  mode: auto # auto | kitty | blocks - inline image renderer
  eager_full_quality: true # download full resolution in the background on chat open
  kitty_placement_cap: 16 # max inline images kept on the terminal at once
  max_long_side_px: 800 # cap a rendered image's long side; height also ≤ 2/3 pane
  disk_cache_size: 268435456 # 256 MB of fetched chat media kept between runs

avatars:
  disk_cache_size: 16777216 # 16 MB of people's pictures, budgeted separately

# keybindings: see keybindings.md
```

## State directory

`state_dir` sets where the account's state lives - the Telegram session, the
local database, and the instance lock. It defaults to `$XDG_STATE_HOME/tele`,
falling back to `~/.local/state/tele`.

Only one `tele` instance can use a state directory at a time. A second one exits
immediately and names the process holding it. Two instances shared one session
and one database with nothing arbitrating between them, quietly overwriting each
other's unread counts and sync state, so this is now refused rather than left to
fail quietly later.

The older `telegram.session_file` key still works and keeps the session where it
points, but it is deprecated and will be removed in the next release. If you have
not set it, your existing session and database are moved into the state directory
automatically on first run - nothing is lost and you stay logged in.

## Notifications

A new message reaches you two ways, and `ui.notifications` switches them
separately. `desktop` hands the notification to your operating system's
notification service, where it survives `tele` not being on screen and can be
routed by your own notification daemon. `toast` draws it in a corner of `tele`'s
own window, which no daemon rule can reach. Turn either off, or both.

Neither switch touches the chat list: with both off, a new message still
highlights its row and moves the chat up. That is the message arriving, not an
interruption - the same line Telegram's own mute draws.

`preview` is about what a notification says rather than where it goes, so it
applies to both: off, the desktop notification and the toast alike carry the
sender's name and nothing else. The body is rendered once and handed to each
unchanged, which is what keeps the two from ever disagreeing about the same
message.

`ui.toasts` places the toasts themselves. Errors, warnings and confirmations go
to `error_zone`; notifications go to `notify_zone`. Both take `bottom-right` or
`top-right` and may name the same corner, in which case they stack together. The
bottom-left corner is not offered: it is kept for the key-press overlay.

## Photos

`photos.mode` picks the renderer. `auto` uses the Kitty graphics protocol in the
terminals known to place images through Unicode placeholders - kitty, Ghostty,
and iTerm2 from 3.7.0 - and draws ANSI block art everywhere else, including
inside tmux and screen, which do not pass the protocol through. `kitty` and
`blocks` force one renderer for a terminal the heuristic does not know: forcing
`kitty` on one that ignores placeholder cells leaves a photo as blank space
rather than block art. The value is read once at startup.

`kitty_placement_cap` bounds how many Kitty image placements are live on the
terminal simultaneously. Only on-screen images (plus a few recently
scrolled-past) are transmitted; older ones are evicted. Transmitting an entire
heavy chat at once can exceed the terminal's image limit and corrupt placements
(shrunken or shifted photos) - lower the cap if you still see that.

`max_long_side_px` caps a rendered inline image's long side in pixels (mirrors
the desktop clients' fixed media size). The height is additionally bounded to 2/3
of the chat pane so a tall photo never dominates the view. Raise it for larger
inline images, lower it for more compact ones.

## Avatars

`avatars.disk_cache_size` is a second budget, deliberately not part of
`photos.disk_cache_size`. An avatar is around a hundred kilobytes and is asked
for again every time you open that person's profile, while chat media is
unbounded and looked at once - sharing one budget would let a scrolling session
evict every face you have. Either key set to `0` means "keep nothing between
runs": that cache moves into a temp directory and is deleted on exit.

## See also

- [Themes](themes.md) - the `ui.theme` slots, writing your own, every token
- [Keybindings](keybindings.md#configurable-actions) - the `keybindings:` section
- [Media](media.md) - what the photo, voice and video settings affect
