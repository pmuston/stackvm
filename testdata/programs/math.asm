; Math functions example
; Calculate sqrt(3^2 + 4^2) = 5 (Pythagorean theorem)

START:
    PUSH 3
    DUP
    MUL             ; 3^2 = 9

    PUSH 4
    DUP
    MUL             ; 4^2 = 16

    ADD             ; 9 + 16 = 25
    SQRT            ; sqrt(25) = 5

    HALT
