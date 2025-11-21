; array_sum.asm
; Sum an array of numbers stored in memory
; Array: [10, 20, 30, 40, 50] at memory[0..4]
; Expected result: 150

    ; Initialize array
    PUSH 10
    STORE 0

    PUSH 20
    STORE 1

    PUSH 30
    STORE 2

    PUSH 40
    STORE 3

    PUSH 50
    STORE 4

    ; Initialize sum and counter
    PUSHI 0     ; sum = 0
    PUSHI 0     ; i = 0

LOOP:
    ; Check if i >= 5 (array length)
    DUP
    PUSHI 5
    GE
    JMPNZ DONE

    ; Load array[i] and add to sum
    DUP         ; Copy i
    LOADD       ; Load array[i]

    ; Add to sum (which is 2nd on stack)
    SWAP        ; [sum, value, i]
    SWAP        ; [sum, i, value]
    ROT         ; [i, value, sum]
    ADD         ; [i, sum+value]
    SWAP        ; [sum, i]

    ; Increment i
    INC
    JMP LOOP

DONE:
    POP         ; Remove counter
    HALT        ; sum is on stack (150)
