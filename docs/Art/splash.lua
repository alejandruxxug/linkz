-- Builds docs/Art/splash.aseprite and the PNG frames/previews for the Linkz
-- splash. Run through the aseprite MCP: dofile(this file).
local ROOT = "/Users/urrexx/Documents/Uni/4th Semester/Langs&Compilers/linkz/"
local ART = ROOT .. "docs/Art/"
local PREV = ART .. "Previews/"

local hex = {
  ink = "1b1a17", night = "2b2a24", ["panel-dark"] = "3a3830", panel = "46443a",
  shadow = "6e6a5a", stone = "a39e88", light = "d6cfb2", paper = "f2ecd2",
  moss = "2f6b3a", leaf = "3fa34d", lime = "8fd16a", sun = "e0a526",
  butter = "f2d06b", rust = "c8642c", ember = "d94a3a", sea = "2d6a8f",
}
local pc = {}
for k, h in pairs(hex) do
  pc[k] = app.pixelColor.rgba(tonumber(h:sub(1, 2), 16), tonumber(h:sub(3, 4), 16), tonumber(h:sub(5, 6), 16), 255)
end

local function img(w, h) return Image(w, h, ColorMode.RGB) end
local function put(im, x, y, name)
  if x >= 0 and y >= 0 and x < im.width and y < im.height then im:drawPixel(x, y, pc[name]) end
end

-- A mask is a set keyed by y*1000+x.
local function key(x, y) return y * 1000 + x end
local function rect(m, x, y, w, h)
  for yy = y, y + h - 1 do for xx = x, x + w - 1 do m[key(xx, yy)] = true end end
end
local function has(m, x, y) return m[key(x, y)] == true end

-- Adds an 8-neighbour ink outline around everything already drawn.
local function outline(im)
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
  for _, p in ipairs(add) do put(im, p[1], p[2], "ink") end
end

---------------------------------------------------------------- logo
local function logo()
  local b = {}
  -- L
  rect(b, 0, 0, 6, 28); rect(b, 0, 22, 20, 6)
  -- i
  rect(b, 24, 0, 6, 6); rect(b, 24, 10, 6, 18)
  -- n
  rect(b, 34, 10, 6, 18); rect(b, 34, 10, 16, 6); rect(b, 46, 12, 6, 16); rect(b, 50, 11, 1, 1)
  -- k
  rect(b, 56, 0, 6, 28)
  for y = 10, 21 do rect(b, 61 + (21 - y), y, 6, 1) end
  for y = 18, 27 do rect(b, 62 + (y - 18), y, 6, 1) end
  -- z
  rect(b, 82, 10, 19, 6); rect(b, 82, 22, 19, 6)
  for y = 16, 21 do rect(b, 94 - (y - 16) * 2, y, 7, 1) end

  local W, H = 101 + 3, 28 + 3
  local im = img(W, H)
  -- shadow first, 1 px down-right, then the letters over it
  for k in pairs(b) do
    local x, y = k % 1000, k // 1000
    put(im, x + 2, y + 2, "moss")
  end
  for k in pairs(b) do
    local x, y = k % 1000, k // 1000
    local top = not has(b, x, y - 1)
    put(im, x + 1, y + 1, top and "lime" or "leaf")
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
local function blit(dst, src, ox, oy)
  for y = 0, src.height - 1 do for x = 0, src.width - 1 do
    local p = src:getPixel(x, y)
    if app.pixelColor.rgbaA(p) > 0 then
      local X, Y = ox + x, oy + y
      if X >= 0 and Y >= 0 and X < dst.width and Y < dst.height then dst:drawPixel(X, Y, p) end
    end
  end end
end

local BW, BH = 200, 80
local L = logo()
local LX, LY = (BW - L.width) // 2, 14
local TAG = "powered by lizard-script"
local TX, TY = (BW - #TAG * 7) // 2, 56
local LZY = LY + L.height - 9  -- back behind the logo, belly and legs below it
local LZX = 120                 -- mid-crawl, head poking out past the z

L:saveAs(PREV .. "logo.png")
scaled(L, 3, true):saveAs(PREV .. "logo-3x.png")

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

local banner = img(BW, BH)
blit(banner, frames[2], LZX, LZY)
blit(banner, L, LX, LY)
blit(banner, tag, TX, TY)
scaled(banner, 3, true):saveAs(PREV .. "banner-3x.png")
banner:saveAs(PREV .. "banner.png")

-- The Aseprite source: banner-sized, one layer each, the walk as a tag.
local spr = Sprite(BW, BH, ColorMode.RGB)
local lz = spr.layers[1]; lz.name = "lizard"
local lg = spr:newLayer(); lg.name = "logo"
local tg = spr:newLayer(); tg.name = "tagline"
for f = 2, 6 do spr:newEmptyFrame() end
for f = 1, 6 do
  spr.frames[f].duration = 1 / 12
  spr:newCel(lz, f, frames[f - 1], Point(LZX, LZY))
  spr:newCel(lg, f, L, Point(LX, LY))
  spr:newCel(tg, f, tag, Point(TX, TY))
end
local tagw = spr:newTag(1, 6); tagw.name = "walk"
local s1 = spr:newSlice(Rectangle(LX, LY, L.width, L.height)); s1.name = "logo"
local s2 = spr:newSlice(Rectangle(LZX, LZY, LW, LH)); s2.name = "lizard"
spr:saveAs(ART .. "splash.aseprite")
print("logo " .. L.width .. "x" .. L.height .. " at " .. LX .. "," .. LY .. "; lizard at " .. LZX .. "," .. LZY)
