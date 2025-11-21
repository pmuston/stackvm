; Loop example
; Count from 1 to 5

START:
    PUSHI 1         ; Counter starts at 1

LOOP:
    DUP             ; Duplicate counter
    PUSHI 5         ; Compare against 5
    GT              ; counter > 5?
    JMPNZ END       ; If yes, exit loop

    INC             ; Increment counter
    JMP LOOP        ; Repeat

END:
    HALT
