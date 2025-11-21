; temperature.asm
; Convert Celsius to Fahrenheit
; Formula: F = (C × 9/5) + 32
; Example: 25°C = 77°F

    PUSH 25     ; Celsius temperature

    ; Multiply by 9
    PUSH 9
    MUL         ; 25 * 9 = 225

    ; Divide by 5
    PUSH 5
    DIV         ; 225 / 5 = 45

    ; Add 32
    PUSH 32
    ADD         ; 45 + 32 = 77°F

    HALT        ; Result: 77
