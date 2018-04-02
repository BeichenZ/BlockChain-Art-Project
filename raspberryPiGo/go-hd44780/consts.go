package hd44780

import "time"

const (
	// Timing constants
	ePulse = 1 * time.Microsecond
	eDelay = 70 * time.Microsecond

	// Some defaults
	lcdWidth = 16 // Maximum characters per line
	lcdChr   = true
	lcdCmd   = false

	lcdLine1 = 0x80 // LCD RAM address for the 1st line
	lcdLine2 = 0xC0 // LCD RAM address for the 2nd line
)
