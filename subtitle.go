package srt

import "time"

type Subtitle struct {
	Number   int
	Start    time.Duration
	End      time.Duration
	Text     string
}

/*
Srt file format:

1
00:00:00,000 --> 00:00:00,000
Blah blah


**repeat**

row 1 is a sequential number starting at 1
row 2 is a timscode, formatted as hours:minutes:seconds,milliseconds
row 3 to an empty row is the content

The content may contain formatting>
Bold: <b> </b> or {b} {/b}
Italic: <i> </i> or {i} {/i}
Underline: <u> </u> or {u} {/u}
Font color: <font color="name or #code"> </font> (HTML colors)

Row 2 may contain DVD rectangle positioning and styling, for ex>

00:00:10,500 --> 00:00:13,000  X1:63 X2:223 Y1:43 Y2:58

or

00:00:15,000 --> 00:00:18,000  X1:53 X2:303 Y1:438 Y2:453

*/
