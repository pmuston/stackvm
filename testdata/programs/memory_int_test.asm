; memory_int_test.asm
; Test program with integer values

    ; Store integers using PUSHI
    PUSHI 42
    STORE 0

    PUSHI 100
    STORE 1

    ; Calculate and store
    PUSHI 10
    PUSHI 5
    ADD
    STORE 2     ; 15

    HALT
