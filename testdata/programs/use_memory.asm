; use_memory.asm
; Test program that uses pre-initialized memory values
; Expects memory[0] and memory[1] to be set before execution

    ; Load values from memory
    LOAD 0      ; Load memory[0]
    LOAD 1      ; Load memory[1]

    ; Add them
    ADD

    ; Store result in memory[2]
    STORE 2

    HALT
