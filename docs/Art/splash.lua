-- Builds docs/Art/splash.aseprite and the PNG frames/previews for the Linkz
-- splash. Run through the aseprite MCP: dofile(this file).
-- The repo root: set LINKZ_ROOT, or run from the repo folder.
local ROOT = (os.getenv("LINKZ_ROOT") or ".") .. "/"
local ART = ROOT .. "docs/Art/"
local PREV = ART .. "Previews/"

local hex = {
  ink = "1b1a17", night = "2b2a24", ["panel-dark"] = "3a3830", panel = "46443a",
  shadow = "6e6a5a", stone = "a39e88", light = "d6cfb2", paper = "f2ecd2",
  moss = "2f6b3a", leaf = "3fa34d", lime = "8fd16a", sun = "e0a526",
  butter = "f2d06b", rust = "c8642c", ember = "d94a3a", sea = "2d6a8f",
  -- "ink" letters (i, n, k): wet blue-violet ink
  inkfill = "3b40a8", inkdark = "1c1f5c", inkedge = "6870e0", inksheen = "a8b0ff", inkglint = "eef0ff",
  -- vines: bluish greens so they stay visible on the green and the ink letters
  -- stargazer lilies (docs/Art/References/lilies-stargazer.png)
  lrim = "fbf7f4", lpink = "f4a6c8", lrose = "e8649e", lline = "c0267a", lspeck = "7e1240",
  lthroat = "c8d84a", ledge = "b77aa6", lthroat2 = "95b83a", lback = "ead2e3", lout = "2a1424",
  lfil = "8fc05a", lanther = "6e2a12", lanther2 = "c8642c",
  bud = "a8467e", budhi = "d778aa", budtip = "e9a4c8", budbase = "6aa84a", spark = "ffffff", spark2 = "fff1a8",
  stem = "2f9470", stemdark = "0f3a2f", young = "c8f08a", leafv = "3fc28e", leafhi = "b4f2d2",
}
local pc = {}
for k, h in pairs(hex) do
  pc[k] = app.pixelColor.rgba(tonumber(h:sub(1, 2), 16), tonumber(h:sub(3, 4), 16), tonumber(h:sub(5, 6), 16), 255)
end

local function img(w, h) return Image(w, h, ColorMode.RGB) end
local function put(im, x, y, name)
  if x >= 0 and y >= 0 and x < im.width and y < im.height then im:drawPixel(x, y, pc[name]) end
end

local function blit(dst, src, ox, oy)
  for y = 0, src.height - 1 do for x = 0, src.width - 1 do
    local p = src:getPixel(x, y)
    if app.pixelColor.rgbaA(p) > 0 then
      local X, Y = ox + x, oy + y
      if X >= 0 and Y >= 0 and X < dst.width and Y < dst.height then dst:drawPixel(X, Y, p) end
    end
  end end
end

-- A mask is a set keyed by y*1000+x.
local function key(x, y) return y * 1000 + x end
local function rect(m, x, y, w, h)
  for yy = y, y + h - 1 do for xx = x, x + w - 1 do m[key(xx, yy)] = true end end
end
local function has(m, x, y) return m[key(x, y)] == true end

-- Adds an 8-neighbour ink outline around everything already drawn.
local function outline(im, color)
  local add = {}
  for y = 0, im.height - 1 do
    for x = 0, im.width - 1 do
      if app.pixelColor.rgbaA(im:getPixel(x, y)) == 0 then
        for dy = -1, 1 do for dx = -1, 1 do
          local nx, ny = x + dx, y + dy
          if nx >= 0 and ny >= 0 and nx < im.width and ny < im.height
            and app.pixelColor.rgbaA(im:getPixel(nx, ny)) > 0 then add[#add + 1] = { x, y } end
        end end
      end
    end
  end
  for _, p in ipairs(add) do put(im, p[1], p[2], color or "ink") end
end

---------------------------------------------------------------- logo
local function logo()
  local g, k = {}, {} -- green letters (L, z) and ink letters (i, n, k)
  -- L
  rect(g, 0, 0, 6, 28); rect(g, 0, 22, 20, 6)
  -- i
  rect(k, 24, 0, 6, 6); rect(k, 24, 10, 6, 18)
  -- n
  rect(k, 34, 10, 6, 18); rect(k, 34, 10, 16, 6); rect(k, 46, 12, 6, 16); rect(k, 50, 11, 1, 1)
  -- k, with an ink drop running off the bottom of its stem
  rect(k, 56, 0, 6, 28)
  rect(k, 58, 28, 2, 2); rect(k, 57, 30, 3, 2); rect(k, 58, 32, 1, 1)
  for y = 10, 21 do rect(k, 61 + (21 - y), y, 6, 1) end
  for y = 18, 27 do rect(k, 62 + (y - 18), y, 6, 1) end
  -- z
  rect(g, 82, 10, 19, 6); rect(g, 82, 22, 19, 6)
  for y = 16, 21 do rect(g, 94 - (y - 16) * 2, y, 7, 1) end

  local b = {}
  for key_ in pairs(g) do b[key_] = true end
  for key_ in pairs(k) do b[key_] = true end
  local W, H = 101 + 3, 33 + 3
  local im = img(W, H)
  -- shadow first, 1 px down-right, then the letters over it
  for key_ in pairs(b) do
    local x, y = key_ % 1000, key_ // 1000
    put(im, x + 2, y + 2, g[key_] and "moss" or "inkdark")
  end
  for key_ in pairs(g) do
    local x, y = key_ % 1000, key_ // 1000
    put(im, x + 1, y + 1, (not has(b, x, y - 1)) and "lime" or "leaf")
  end
  for key_ in pairs(k) do
    local x, y = key_ % 1000, key_ // 1000
    local top, left = not has(b, x, y - 1), not has(b, x - 1, y)
    local c = "inkfill"
    if top then c = "inksheen" elseif left then c = "inkedge" end
    put(im, x + 1, y + 1, c)
  end
  -- wet-ink glints: a short bright streak just inside each top-left corner
  for key_ in pairs(k) do
    local x, y = key_ % 1000, key_ // 1000
    if not has(b, x, y - 1) and not has(b, x - 1, y) and has(b, x + 1, y + 2) and has(b, x + 2, y) then
      put(im, x + 2, y + 1, "inkglint")
      put(im, x + 2, y + 2, "inkglint")
    end
  end
  outline(im)
  return im
end

---------------------------------------------------------------- lizard
local LW, LH = 48, 16
-- stance moves the foot back 2 px a frame (so the body advances 2 px/frame),
-- then two lifted frames swing it forward.
local DX = { 3, 1, -1, -3, -1, 2 }
local LIFT = { 0, 0, 0, 0, 2, 1 }

local function column(x, f)
  if x >= 1 and x <= 16 then
    local t = (x - 1) / 15
    local h = 0.5 + 1.5 * t
    local wave = math.sin(2 * math.pi * f / 6 + (1 - t) * 2.2)
    local c = 10.5 - 2.5 * t + 1.6 * (1 - t) * wave
    return math.floor(c - h + 0.5), math.floor(c + h - 0.5)
  elseif x <= 18 then return 7, 9
  elseif x <= 34 then return 6, 9
  elseif x <= 36 then return 7, 9
  end
  local tops = { [37] = 5, [38] = 4, [39] = 4, [40] = 4, [41] = 4, [42] = 4, [43] = 5, [44] = 5, [45] = 6, [46] = 7 }
  local bots = { [46] = 8 }
  if tops[x] then return tops[x], bots[x] or 9 end
  return nil
end

-- A leg: thigh out of the body, shin down to the foot, toes forward.
local function leg(im, hx, phase, color)
  local dx, lift = DX[phase + 1], LIFT[phase + 1]
  local fy = 14 - lift
  local fx = hx + dx
  for y = 10, fy do
    local t = (fy == 10) and 1 or (y - 10) / (fy - 10)
    local x = math.floor(hx + (fx - hx) * t + 0.5)
    put(im, x, y, color); put(im, x + 1, y, color)
  end
  put(im, fx + 2, fy, color)
end

local function lizard(f)
  local A, B = f % 6, (f + 3) % 6
  -- far legs get their own outline and sit behind the body
  local far = img(LW, LH)
  leg(far, 23, B, "moss")
  leg(far, 35, A, "moss")
  outline(far)
  local im = img(LW, LH)
  for x = 1, 46 do
    local top, bot = column(x, f)
    if top then
      for y = top, bot do
        local c = "leaf"
        if y == top and x >= 3 then c = "lime" elseif y == bot and bot > top then c = "moss" end
        put(im, x, y, c)
      end
    end
  end
  for _, p in ipairs({ { 18, 8 }, { 22, 7 }, { 26, 8 }, { 30, 7 }, { 34, 8 } }) do put(im, p[1], p[2], "moss") end
  leg(im, 19, A, "leaf")
  leg(im, 31, B, "leaf")
  put(im, 40, 5, "sun"); put(im, 41, 5, "ink")
  put(im, 44, 8, "ink"); put(im, 45, 8, "ink")
  outline(im)
  local out = img(LW, LH)
  out:drawImage(far, Point(0, 0))
  out:drawImage(im, Point(0, 0))
  return out
end

---------------------------------------------------------------- vines
-- Vines grow over the logo in STAGES stages, then flowers bloom in BLOOM
-- frames. Coordinates are logo-image coordinates; the vine canvas adds a
-- margin (VX, VY) so tips can curl above the letters.
local STAGES, BLOOM = 24, 10
local VX, VY = 48, 14       -- the logo's place on the vine canvas
local VW, VH = 200, 80      -- the vine canvas is the whole banner

-- A smooth 1 px path through control points (Catmull-Rom), as an ordered
-- list of 8-connected pixels with the L-corners removed.
local function path(pts)
  local out = {}
  local function add(x, y)
    local l = out[#out]
    if l and l[1] == x and l[2] == y then return end
    if l and (math.abs(l[1] - x) > 1 or math.abs(l[2] - y) > 1) then
      local n = math.max(math.abs(l[1] - x), math.abs(l[2] - y))
      for i = 1, n - 1 do
        local px = math.floor(l[1] + (x - l[1]) * i / n + 0.5)
        local py = math.floor(l[2] + (y - l[2]) * i / n + 0.5)
        local q = out[#out]
        if q[1] ~= px or q[2] ~= py then out[#out + 1] = { px, py } end
      end
    end
    out[#out + 1] = { x, y }
  end
  local n = #pts
  for i = 1, n - 1 do
    local p0, p1, p2, p3 = pts[math.max(1, i - 1)], pts[i], pts[i + 1], pts[math.min(n, i + 2)]
    for step = 0, 39 do
      local t = step / 40
      local t2, t3 = t * t, t * t * t
      local function cr(a0, a1, a2, a3)
        return 0.5 * ((2 * a1) + (-a0 + a2) * t + (2 * a0 - 5 * a1 + 4 * a2 - a3) * t2 + (-a0 + 3 * a1 - 3 * a2 + a3) * t3)
      end
      add(math.floor(cr(p0[1], p1[1], p2[1], p3[1]) + 0.5), math.floor(cr(p0[2], p1[2], p2[2], p3[2]) + 0.5))
    end
  end
  add(pts[n][1], pts[n][2])
  -- pixel-perfect: drop a pixel whose neighbours already touch diagonally
  local clean = { out[1] }
  for i = 2, #out - 1 do
    local a, c = clean[#clean], out[i + 1]
    if not (math.abs(a[1] - c[1]) == 1 and math.abs(a[2] - c[2]) == 1) then clean[#clean + 1] = out[i] end
  end
  clean[#clean + 1] = out[#out]
  return clean
end

-- Branches: control points, the parent they sprout from, the stage they
-- finish, and leaves {fraction along the branch, direction}.
local branches = {
  { name = "L", pts = { { 2, 33 }, { 0, 29 }, { -1, 25 }, { 7, 19 }, { 0, 13 }, { 7, 7 }, { 3, 2 }, { 5, -2 }, { 8, -3 }, { 10, -1 }, { 9, 0 } },
    start = 0, done = 12, leaves = { { 0.04, "ur" }, { 0.18, "dl" }, { 0.42, "ur" }, { 0.66, "ul" }, { 0.86, "ur" } } },
  { name = "runner", parent = "L", pts = { { 2, 33 }, { 12, 34 }, { 22, 32 }, { 32, 34 }, { 44, 32 }, { 56, 35 }, { 68, 32 }, { 80, 34 }, { 92, 32 }, { 102, 34 }, { 107, 30 }, { 106, 26 }, { 104, 27 }, { 105, 28 } },
    done = 18, leaves = { { 0.1, "dr" }, { 0.24, "ul" }, { 0.38, "dr" }, { 0.55, "dl" }, { 0.7, "dr" }, { 0.84, "dl" } } },
  { name = "west", parent = "L", pts = { { 1, 31 }, { -5, 33 }, { -11, 30 }, { -15, 24 }, { -18, 20 }, { -17, 16 } },
    done = 16, leaves = { { 0.3, "dl" }, { 0.55, "ur" }, { 0.78, "dl" } }, lily = { -19, 5 } },
  { name = "i", parent = "runner", pts = { { 29, 33 }, { 31, 27 }, { 23, 21 }, { 31, 14 }, { 30, 9 }, { 23, 7 }, { 22, 2 }, { 25, -2 }, { 28, -3 }, { 29, -1 } },
    done = 18, leaves = { { 0.3, "ul" }, { 0.62, "ur" } } },
  { name = "n", parent = "runner", pts = { { 50, 33 }, { 54, 26 }, { 46, 19 }, { 53, 12 }, { 50, 8 }, { 44, 6 }, { 41, 2 } },
    done = 19, leaves = { { 0.35, "dr" }, { 0.7, "ur" } }, bud = true },
  { name = "k", parent = "runner", pts = { { 62, 33 }, { 64, 26 }, { 55, 20 }, { 63, 13 }, { 55, 7 }, { 59, 0 }, { 62, -3 }, { 64, -2 }, { 63, -1 } },
    done = 20, leaves = { { 0.25, "ul" }, { 0.55, "ur" }, { 0.8, "ul" } } },
  { name = "z", parent = "runner", pts = { { 101, 33 }, { 104, 25 }, { 103, 18 }, { 101, 12 }, { 97, 8 }, { 91, 8 }, { 87, 7 }, { 85, 5 }, { 87, 4 } },
    done = 21, leaves = { { 0.3, "dr" }, { 0.6, "ur" } } },
  { name = "east", parent = "runner", pts = { { 100, 34 }, { 108, 35 }, { 114, 31 }, { 118, 25 }, { 121, 20 }, { 122, 17 } },
    done = 22, leaves = { { 0.3, "ul" }, { 0.55, "dr" }, { 0.8, "ul" } }, lily = { 125, 6 } },
}
local byname = {}
for _, b in ipairs(branches) do
  b.px = path(b.pts)
  byname[b.name] = b
end

-- How many pixels of a branch exist at a stage. A child starts the stage
-- after its parent's tip passes the point it grows from.
local function nearest(px, p)
  local best, bi = 1e9, 1
  for i, q in ipairs(px) do
    local d = (q[1] - p[1]) ^ 2 + (q[2] - p[2]) ^ 2
    if d < best then best, bi = d, i end
  end
  return bi
end
local function shown(b, s)
  if b.start == nil then
    local par = byname[b.parent]
    local idx = nearest(par.px, b.pts[1])
    b.start = STAGES
    for t = 0, STAGES - 1 do
      if shown(par, t) >= idx then b.start = t + 1; break end
    end
  end
  if s < b.start then return 0 end
  local f = (s - b.start + 1) / (b.done - b.start + 1)
  if b.start == 0 then f = ((s + 1) / (b.done + 1)) ^ 1.3 end -- the sprout starts slow
  local least = (b.start == 0) and 4 or 1 -- stage 0 is already a little sprout
  return math.min(#b.px, math.max(least, math.floor(#b.px * f + 0.5)))
end

local LEAF = {
  ur = { { 1, -1, "leafv" }, { 2, -1, "leafv" }, { 2, -2, "leafhi" }, { 3, -2, "leafv" } },
  ul = { { -1, -1, "leafv" }, { -2, -1, "leafv" }, { -2, -2, "leafhi" }, { -3, -2, "leafv" } },
  dr = { { 1, 1, "leafv" }, { 2, 1, "leafhi" }, { 2, 2, "leafv" }, { 3, 2, "leafv" } },
  dl = { { -1, 1, "leafv" }, { -2, 1, "leafhi" }, { -2, 2, "leafv" }, { -3, 2, "leafv" } },
}

local draw_bud -- defined with the lilies below

-- One vine stage as an image on the vine canvas.
local function vine_stage(s)
  local cells = {} -- [key] = color name, drawn later in two passes
  local order = {}
  local function set(x, y, c)
    local k = (y + 100) * 1000 + (x + 100)
    if cells[k] == nil then order[#order + 1] = k end
    cells[k] = c
  end
  for _, b in ipairs(branches) do
    local n = shown(b, s)
    for i = 1, n do
      local p = b.px[i]
      set(p[1], p[2], (n < #b.px and i > n - 2) and "young" or "stem")
    end
    for _, lf in ipairs(b.leaves) do
      local i = math.max(1, math.floor(#b.px * lf[1]))
      local age = n - i -- pixels grown past the leaf
      if age >= 2 then
        local p = b.px[i]
        local parts = LEAF[lf[2]]
        local m = (age >= 6) and #parts or ((age >= 4) and 3 or 1) -- bud, half, full
        for j = 1, m do set(p[1] + parts[j][1], p[2] + parts[j][2], parts[j][3]) end
      end
    end
  end
  local im = img(VW, VH)
  local budded = {}
  for _, b in ipairs(branches) do
    if b.bud and s >= 21 then budded[#budded + 1] = { b, s >= 22 } end
  end
  -- a dark edge 1 px down-right of every vine pixel keeps it apart from the letters
  for _, k in ipairs(order) do
    local x, y = k % 1000 - 100, k // 1000 - 100
    local below = (y + 1 + 100) * 1000 + (x + 1 + 100)
    if cells[below] == nil then put(im, x + 1 + VX, y + 1 + VY, "stemdark") end
  end
  for _, k in ipairs(order) do
    local x, y = k % 1000 - 100, k // 1000 - 100
    put(im, x + VX, y + VY, cells[k])
  end
  for _, bb in ipairs(budded) do
    local tip = bb[1].px[#bb[1].px]
    draw_bud(im, tip[1] + VX, tip[2] + VY, bb[2] and 9 or 6, bb[2] and 4 or 3, -math.pi / 2 - 0.5, false)
  end
  return im
end

---------------------------------------------------------------- lilies
-- Stargazer lilies: six wide pointed petals, white rim, pink body, a magenta
-- line down each, crimson speckles, a yellow-green star throat, green
-- stamens with dark anthers. Drawn into a scratch image, outlined, then
-- copied onto the target so the outline never touches the vines.
local function hash(x, y) return ((x * 73856093) ~ (y * 19349663)) % 97 end

-- A long closed bud from (x, y) toward angle a: len px long, wid px wide.
draw_bud = function(dst, x, y, len, wid, a, pinker)
  local im = img(dst.width, dst.height)
  local ux, uy = math.cos(a), math.sin(a)
  for py = math.floor(y - len - 2), math.ceil(y + len + 2) do
    for px = math.floor(x - len - 2), math.ceil(x + len + 2) do
      local dx, dy = px - x, py - y
      local u = dx * ux + dy * uy          -- along the bud
      local v = -dx * uy + dy * ux         -- across it
      local t = u / len
      if t >= 0 and t <= 1 then
        local half = wid / 2 * math.sin(math.pi * math.min(1, t * 1.15) ^ 0.8) + 0.25
        if math.abs(v) <= half then
          local c = pinker and "budhi" or "bud"
          if t < 0.18 then c = "budbase"
          elseif t > 0.82 then c = "budtip"
          elseif v < -half + 1 then c = pinker and "budtip" or "budhi" end
          put(im, px, py, c)
        end
      end
    end
  end
  outline(im, "lout")
  blit(dst, im, 0, 0)
end

-- One lily. open: 0 = just split, 1 = fully open. stamen: 0..1.
local function draw_lily(dst, cx, cy, R, open, up, stamen)
  local im = img(dst.width, dst.height)
  local SQ = 0.86 -- tilted a little toward the viewer
  local W = 0.40  -- petal half-width at its widest, as a part of R
  local petals = {}
  for i = 0, 5 do
    local full = up + math.pi / 6 + i * math.pi / 3
    local d = math.atan(math.sin(full - up), math.cos(full - up))
    petals[#petals + 1] = { a = up + d * open, front = (i % 2 == 0), r = R * (0.55 + 0.45 * open) * ((i % 2 == 0) and 1 or 0.92) }
  end
  -- back petals first, then the front three over them
  local function petal(p, pass)
    if p.front ~= pass then return end
    local ux, uy = math.cos(p.a), math.sin(p.a)
    for py = math.floor(cy - R - 2), math.ceil(cy + R + 2) do
      for px = math.floor(cx - R - 2), math.ceil(cx + R + 2) do
        local dx, dy = px - cx, (py - cy) / SQ
        local u = dx * ux + dy * uy
        local v = -dx * uy + dy * ux
        local t = u / p.r
        local curl = 0.12 * open * p.r * t * t -- tips sweep back sideways
        local vv = v - curl
        if t > 0 and t <= 1 then
          local half = W * p.r * math.sin(math.pi * t ^ 0.62) ^ 1.25 + 0.35
          if math.abs(vv) <= half then
            local s = math.abs(vv) / half
            local c
            if s > 0.66 or t > 0.86 then c = "lrim"
            elseif s > 0.36 or t > 0.66 then c = "lpink"
            else c = "lrose" end
            if math.abs(vv) < 0.55 and t > 0.3 and t < 0.82 then c = "lline" end
            if math.abs(vv) < 0.55 and t <= 0.3 then c = "lthroat" end
            if t < 0.14 then c = "lthroat2" end
            if not p.front and c == "lrim" then c = "lback" end
            if t > 0.92 and open > 0.9 then c = "lback" end -- the tip curls away
            -- where a front petal lies over a back one, its edge is shaded
            if p.front and s > 0.7 and app.pixelColor.rgbaA(im:getPixel(px, py)) > 0 then c = "ledge" end
            put(im, px, py, c)
          end
        end
      end
    end
  end
  for _, p in ipairs(petals) do petal(p, false) end
  -- a thin shade between front and back petals so the overlap reads
  for _, p in ipairs(petals) do petal(p, true) end
  -- speckles: two on each front petal, inside the pink
  if open > 0.6 then
    for _, p in ipairs(petals) do
      if p.front then
        for _, sp in ipairs({ { 0.48, 1.3 }, { 0.6, -1.4 } }) do
          local ux, uy = math.cos(p.a), math.sin(p.a)
          local u, v = sp[1] * p.r, sp[2]
          local px = math.floor(cx + u * ux - v * uy + 0.5)
          local py = math.floor(cy + (u * uy + v * ux) * SQ + 0.5)
          put(im, px, py, "lspeck")
        end
      end
    end
  end
  outline(im, "lout")
  -- stamens: three green filaments with dark anthers, poking out of the throat
  if stamen > 0 then
    for i, a in ipairs({ up + 0.35, up - 0.5, up + 1.25 }) do
      local len = R * (0.42 + 0.1 * i) * stamen
      local ex, ey = cx, cy
      for k = 1, math.floor(len) do
        ex = math.floor(cx + math.cos(a) * k + 0.5); ey = math.floor(cy + math.sin(a) * k * SQ + 0.5)
        if k > 3 then put(im, ex, ey, "lfil") end
      end
      if stamen >= 1 then
        put(im, ex, ey, "lanther"); put(im, ex + 1, ey, "lanther2")
      end
    end
  end
  blit(dst, im, 0, 0)
end

local function sparkle(dst, x, y)
  put(dst, x, y, "spark2")
  for _, o in ipairs({ { 1, 0 }, { -1, 0 }, { 0, 1 }, { 0, -1 } }) do put(dst, x + o[1], y + o[2], "spark") end
  put(dst, x + 2, y, "spark2"); put(dst, x - 2, y, "spark2"); put(dst, x, y - 2, "spark2"); put(dst, x, y + 2, "spark2")
end

-- The lilies, in the order they open, two frames apart.
local lilies = {
  { branch = "west", R = 14, up = -math.pi / 2 - 0.35, start = 0 },
  { branch = "east", R = 14, up = -math.pi / 2 + 0.3, start = 2 },
}
local function bloom_frame(f)
  local im = img(VW, VH)
  for _, l in ipairs(lilies) do
    local b = byname[l.branch]
    local tip = b.px[#b.px]
    local tx, ty = tip[1] + VX, tip[2] + VY
    local cx, cy = b.lily[1] + VX, b.lily[2] + VY
    local age = f - l.start
    if age == 0 then draw_bud(im, tx, ty, 9, 4, l.up, false)
    elseif age == 1 then draw_bud(im, tx, ty, 12, 6, l.up, true)
    elseif age >= 2 then
      local open = ({ 0.2, 0.5, 0.8, 1, 1, 1 })[math.min(6, age - 1)]
      local st = ({ 0, 0, 0, 0.5, 1, 1 })[math.min(6, age - 1)]
      -- the flower's heart rises from the bud's base to its place
      local k = math.min(1, (age - 1) / 3)
      draw_lily(im, math.floor(tx + (cx - tx) * k + 0.5), math.floor(ty + (cy - ty) * k + 0.5), l.R, open, l.up, st)
      if age == 7 then sparkle(im, cx + 9, cy - 9) end
    end
  end
  return im
end

---------------------------------------------------------------- font
local glyphs = {}
do
  local fh = io.open(ROOT .. "lib/gfx/font_data.liz")
  for line in fh:lines() do
    local ch, rest = line:match('^  "([^"\\]+)": (%[.*)$')
    if ch and #ch == 1 then
      local runs = {}
      for y, a, b in rest:gmatch("%[(%d+), (%d+), (%d+)%]") do runs[#runs + 1] = { tonumber(y), tonumber(a), tonumber(b) } end
      glyphs[ch] = runs
    end
  end
  fh:close()
end
local function text(im, s, x, y, color)
  for i = 1, #s do
    for _, r in ipairs(glyphs[s:sub(i, i)] or {}) do
      for xx = r[2], r[3] - 1 do put(im, x + (i - 1) * 7 + xx, y + r[1], color) end
    end
  end
end

---------------------------------------------------------------- outputs
local function scaled(src, s, bg)
  local out = img(src.width * s, src.height * s)
  for y = 0, out.height - 1 do
    for x = 0, out.width - 1 do
      local p = src:getPixel(x // s, y // s)
      if app.pixelColor.rgbaA(p) == 0 and bg then
        local cx, cy = (x // s) // 4, (y // s) // 4
        p = ((cx + cy) % 2 == 0) and app.pixelColor.rgba(204, 204, 204, 255) or app.pixelColor.rgba(255, 255, 255, 255)
      end
      out:drawPixel(x, y, p)
    end
  end
  return out
end
local function dark(src, s)
  local out = img(src.width * s, src.height * s)
  out:clear(app.pixelColor.rgba(0x12, 0x12, 0x12, 255))
  for y = 0, out.height - 1 do for x = 0, out.width - 1 do
    local p = src:getPixel(x // s, y // s)
    if app.pixelColor.rgbaA(p) > 0 then out:drawPixel(x, y, p) end
  end end
  return out
end

local BW, BH = 200, 80
local L = logo()
local LX, LY = (BW - L.width) // 2, 14
local TAG = "powered by lizard-script"
local TX, TY = (BW - #TAG * 7) // 2, 56
local LZY = LY + 22             -- back behind the logo, belly and legs below it
local LZX = 120                 -- mid-crawl, head poking out past the z

L:saveAs(PREV .. "logo.png")
scaled(L, 3, true):saveAs(PREV .. "logo-3x.png")
dark(L, 3):saveAs(PREV .. "logo-dark-3x.png")

local frames = {}
local sheet = img(LW * 6 + 5 * 4, LH)
for f = 0, 5 do
  frames[f] = lizard(f)
  frames[f]:saveAs(PREV .. "lizard-" .. f .. ".png")
  blit(sheet, frames[f], f * (LW + 4), 0)
end
scaled(sheet, 4, true):saveAs(PREV .. "lizard-walk-4x.png")

local tag = img(#TAG * 7 + 1, 14)
text(tag, TAG, 1, 1, "ink")
text(tag, TAG, 0, 0, "stone")

local VBX, VBY = LX - VX, LY - VY -- (0, 0): the vine canvas is the banner

local stages, blooms = {}, {}
for st = 0, STAGES - 1 do
  stages[st] = vine_stage(st)
  stages[st]:saveAs(string.format("%svines-%02d.png", PREV, st))
end
for f = 0, BLOOM - 1 do
  blooms[f] = bloom_frame(f)
  blooms[f]:saveAs(string.format("%sbloom-%d.png", PREV, f))
end
for _, b in ipairs(branches) do print(b.name, "starts", b.start, "done", b.done, "px", #b.px) end

local function compose(lizf, lizx, st, bf)
  local im = img(BW, BH)
  blit(im, frames[lizf], lizx, LZY)
  blit(im, L, LX, LY)
  if st then blit(im, stages[st], VBX, VBY) end
  if bf then blit(im, blooms[bf], VBX, VBY) end
  blit(im, tag, TX, TY)
  return im
end

local banner = compose(2, LZX)
scaled(banner, 3, true):saveAs(PREV .. "banner-3x.png")
banner:saveAs(PREV .. "banner.png")
dark(compose(2, LZX, 0), 3):saveAs(PREV .. "banner-vines-0-dark-3x.png")
dark(compose(2, LZX, STAGES // 2), 3):saveAs(PREV .. "banner-vines-50-dark-3x.png")
dark(compose(2, LZX, STAGES - 1), 3):saveAs(PREV .. "banner-vines-100-dark-3x.png")
scaled(compose(2, LZX, STAGES - 1), 3, true):saveAs(PREV .. "banner-vines-100-3x.png")
dark(compose(2, LZX, STAGES - 1, BLOOM - 1), 3):saveAs(PREV .. "banner-final-dark-3x.png")
scaled(compose(2, LZX, STAGES - 1, BLOOM - 1), 3, true):saveAs(PREV .. "banner-final-3x.png")

-- contact sheets: logo + vines per stage (4 across), logo + full vines + bloom per frame (3 across)
local function sheet_of(n, across, draw)
  local cw, ch = VW + 4, VH + 4
  local sh = img(cw * across, ch * math.ceil(n / across))
  for i = 0, n - 1 do
    local cell = img(VW, VH)
    blit(cell, L, VX, VY)
    draw(cell, i)
    blit(sh, cell, (i % across) * cw + 2, (i // across) * ch + 2)
  end
  return sh
end
dark(sheet_of(STAGES, 4, function(c, i) blit(c, stages[i], 0, 0) end), 2):saveAs(PREV .. "vines-sheet-dark-2x.png")
dark(sheet_of(BLOOM, 3, function(c, i) blit(c, stages[STAGES - 1], 0, 0); blit(c, blooms[i], 0, 0) end), 2):saveAs(PREV .. "bloom-sheet-dark-2x.png")

-- The Aseprite source: banner-sized, one layer each. Tags: "walk" (the
-- lizard, 6 frames), "grow" (24 vine stages), "bloom" (9 flower frames).
local spr = Sprite(BW, BH, ColorMode.RGB)
local lz = spr.layers[1]; lz.name = "lizard"
local lg = spr:newLayer(); lg.name = "logo"
local vn = spr:newLayer(); vn.name = "vines"
local bl = spr:newLayer(); bl.name = "bloom"
local tg = spr:newLayer(); tg.name = "tagline"
local total = 6 + STAGES + BLOOM
for f = 2, total do spr:newEmptyFrame() end
for f = 1, total do
  spr.frames[f].duration = 1 / 12
  spr:newCel(lz, f, frames[(f - 1) % 6], Point(LZX, LZY))
  spr:newCel(lg, f, L, Point(LX, LY))
  spr:newCel(tg, f, tag, Point(TX, TY))
  local g = f - 7
  if g >= 0 then spr:newCel(vn, f, stages[math.min(g, STAGES - 1)], Point(VBX, VBY)) end
  local bf = f - 7 - STAGES
  if bf >= 0 then spr:newCel(bl, f, blooms[bf], Point(VBX, VBY)) end
end
local t1 = spr:newTag(1, 6); t1.name = "walk"
local t2 = spr:newTag(7, 6 + STAGES); t2.name = "grow"
local t3 = spr:newTag(7 + STAGES, total); t3.name = "bloom"
local s1 = spr:newSlice(Rectangle(LX, LY, L.width, L.height)); s1.name = "logo"
local s2 = spr:newSlice(Rectangle(LZX, LZY, LW, LH)); s2.name = "lizard"
local s3 = spr:newSlice(Rectangle(VBX, VBY, VW, VH)); s3.name = "vines"
spr:saveAs(ART .. "splash.aseprite")
print("logo " .. L.width .. "x" .. L.height .. " at " .. LX .. "," .. LY .. "; lizard at " .. LZX .. "," .. LZY .. "; vines " .. VW .. "x" .. VH .. " at " .. VBX .. "," .. VBY)
