; memory_test.asm
; Test program to verify memory inspection flag
; Stores various values in memory

    ; Store some integers
    PUSH 42
    STORE 0

    PUSH 100
    STORE 1

    ; Calculate and store a result
    PUSH 10
    PUSH 5
    ADD
    STORE 2     ; 15

    ; Store a float
    PUSH 3.14
    STORE 3

    ; Store result of calculation
    PUSH 5
    PUSH 3
    MUL
    STORE 4     ; 15

    HALT
