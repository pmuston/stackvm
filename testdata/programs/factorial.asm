; factorial.asm
; Calculate 5! = 5 * 4 * 3 * 2 * 1 = 120

    PUSHI 5     ; n = 5
    PUSHI 1     ; result = 1

LOOP:
    OVER        ; Copy n
    PUSHI 1
    LE          ; n <= 1?
    JMPNZ DONE

    OVER        ; Copy n
    MUL         ; result *= n

    SWAP        ; Get n to top
    DEC         ; n--
    SWAP        ; Put result back

    JMP LOOP

DONE:
    SWAP        ; Get result to top
    POP         ; Remove n
    HALT        ; Result: 120
