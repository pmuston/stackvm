; fibonacci_simple.asm
; Calculate Fibonacci numbers using memory
; Stores first 10 Fibonacci numbers in memory[0..9]

    ; Initialize first two numbers
    PUSHI 0
    STORE 0     ; fib(0) = 0

    PUSHI 1
    STORE 1     ; fib(1) = 1

    ; Calculate fib(2) through fib(9)
    PUSHI 2     ; i = 2

LOOP:
    DUP
    PUSHI 10
    GE          ; i >= 10?
    JMPNZ DONE

    ; fib(i) = fib(i-1) + fib(i-2)
    DUP         ; Copy i
    DEC         ; i-1
    LOADD       ; Load fib(i-1)

    OVER        ; Copy i
    DEC
    DEC         ; i-2
    LOADD       ; Load fib(i-2)

    ADD         ; fib(i-1) + fib(i-2)

    OVER        ; Copy i
    SWAP        ; Get result on top
    STORED      ; Store at memory[i]

    INC         ; i++
    JMP LOOP

DONE:
    POP         ; Remove counter

    ; Load final result (fib(9) = 34)
    PUSHI 9
    LOADD
    HALT
