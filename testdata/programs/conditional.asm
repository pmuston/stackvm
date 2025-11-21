; Conditional jump example
; If value > 10, push 1, else push 0

START:
    PUSH 15         ; Test value
    PUSH 10         ; Compare against 10
    GT              ; value > 10?
    JMPZ BELOW      ; Jump if not greater
    PUSH 1          ; Result = 1
    JMP DONE
BELOW:
    PUSH 0          ; Result = 0
DONE:
    HALT
