import AppKit

guard CommandLine.arguments.count == 2 else {
    fatalError("usage: generate_app_icon.swift <output.png>")
}

let image = NSImage(size: NSSize(width: 1024, height: 1024))
image.lockFocus()

NSColor(calibratedRed: 0.11, green: 0.12, blue: 0.16, alpha: 1).setFill()
NSBezierPath(roundedRect: NSRect(x: 72, y: 72, width: 880, height: 880), xRadius: 210, yRadius: 210).fill()

let red = NSColor(calibratedRed: 0.94, green: 0.15, blue: 0.18, alpha: 1)
red.setFill()
NSBezierPath(roundedRect: NSRect(x: 250, y: 300, width: 524, height: 430), xRadius: 210, yRadius: 210).fill()

func strokeLine(from start: NSPoint, to end: NSPoint, width: CGFloat) {
    let line = NSBezierPath()
    line.move(to: start)
    line.line(to: end)
    line.lineWidth = width
    line.lineCapStyle = .round
    line.stroke()
}

red.setStroke()
strokeLine(from: NSPoint(x: 390, y: 705), to: NSPoint(x: 330, y: 835), width: 28)
strokeLine(from: NSPoint(x: 634, y: 705), to: NSPoint(x: 694, y: 835), width: 28)

NSColor.white.setFill()
for x in [390.0, 570.0] {
    NSBezierPath(ovalIn: NSRect(x: x, y: 535, width: 78, height: 92)).fill()
}
NSColor(calibratedWhite: 0.08, alpha: 1).setFill()
for x in [414.0, 594.0] {
    NSBezierPath(ovalIn: NSRect(x: x, y: 555, width: 34, height: 48)).fill()
}

let smile = NSBezierPath()
smile.move(to: NSPoint(x: 420, y: 455))
smile.curve(to: NSPoint(x: 604, y: 455), controlPoint1: NSPoint(x: 465, y: 385), controlPoint2: NSPoint(x: 559, y: 385))
smile.lineWidth = 24
smile.lineCapStyle = .round
NSColor.white.setStroke()
smile.stroke()

red.setStroke()
for rect in [NSRect(x: 125, y: 230, width: 260, height: 240), NSRect(x: 639, y: 230, width: 260, height: 240)] {
    let claw = NSBezierPath(ovalIn: rect)
    claw.lineWidth = 65
    claw.stroke()
}

image.unlockFocus()

guard
    let tiff = image.tiffRepresentation,
    let bitmap = NSBitmapImageRep(data: tiff),
    let png = bitmap.representation(using: .png, properties: [:])
else {
    fatalError("failed to render icon")
}

try png.write(to: URL(filePath: CommandLine.arguments[1]))
