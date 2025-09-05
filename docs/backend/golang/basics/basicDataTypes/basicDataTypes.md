# Go Basic Data Types Study Notes

---

## 1. Integer Types

Go supports multiple bases for integers: decimal, hexadecimal, octal, and binary.

```go
// Hexadecimal
var hex1 uint8 = 0xF   // Uppercase F
var hex2 uint8 = 0xf   // Lowercase f

// Octal
var oct1 uint8 = 017    // Traditional octal
var oct2 uint8 = 0o17   // Go 1.13+ octal notation
var oct3 uint8 = 0o17   // Can be reused

// Binary
var bin1 uint8 = 0b1111 // Go 1.13+ binary
var bin2 uint8 = 0b1111

// Decimal
var dec uint8 = 15

fmt.Println("Integer representations:", hex1, hex2, oct1, oct2, oct3, bin1, bin2, dec)
```

**Notes:**
- Go supports various integer types: `uint8`, `uint16`, `uint32`, `uint64`, etc.
- Binary, octal, and hexadecimal are often used in low-level or hardware-related development.

---

## 2. Float Types

Go supports `float32` and `float64`.

```go
var float1 float32 = 10      // Explicit float32
float2 := 10.0               // Default is float64
fmt.Println("float1 == float2? ", float1 == float32(float2))
```

**Notes:**
- The default float type is `float64`.
- Be careful with type conversion when comparing floats.

---

## 3. Complex Types

Go supports `complex64` and `complex128`, and provides the `complex()` function.

```go
var c1 complex64 = 1 + 2i      // Explicit complex64
c2 := 1 + 2i                   // Default is complex128
c3 := complex(1, 2)            // Using complex() function

fmt.Println("c1 == complex64(c2)? ", c1 == complex64(c2))
fmt.Println("c3 = ", c3)

x := real(c3)  // Real part
y := imag(c3)  // Imaginary part
fmt.Println("Real and imaginary parts of c3:", x, y)
```

**Notes:**
- `real()` gets the real part, `imag()` gets the imaginary part.
- Complex numbers are often used in signal processing and scientific computing.

---

## 4. Strings and Byte Slices

**String to byte slice:**
```go
var str string = "hello"
var bytes []byte = []byte(str)
fmt.Println("String to bytes:", bytes)
```

**Byte slice to string:**
```go
var str1 string = string([]byte{97, 98, 99, 100})
fmt.Println("Bytes to string:", str1)
```

---

## 5. Strings and Unicode (Rune)

**String to rune slice:**
```go
var str2 string = "abc，你好，世界！"
var runes []rune = []rune(str2)
fmt.Println("String to rune:", runes)
```

**Rune slice to string:**
```go
var str3 string = string(runes)
fmt.Println("Rune to string:", str3)
```

**Notes:**
- `rune` is an alias for `int32`, used for Unicode characters.
- Use rune slices to correctly handle Chinese or other multi-byte characters.

---

## 6. Raw Strings and Escaped Strings

```go
var s1 string = "Hello\nworld!\n" // Escaped string
var s2 string = `Hello
    world!
    ` // Raw string
fmt.Println("s1 == s2? ", s1 == s2)
```

**Notes:**
- Escaped strings support `\n`, `\t`, etc.
- Raw strings use backticks `` ` `` and preserve newlines and spaces.

---

## 7. String Length and Slicing

```go
var str4 string = "Hello, 世界"
var bytes2 []byte = []byte(str4)
var runes2 []rune = []rune(str4)

fmt.Println("len(str4):", len(str4))         // Byte length
fmt.Println("len(bytes2):", len(bytes2))
fmt.Println("len(runes2):", len(runes2))     // Character length

// Slicing
fmt.Println("str4[0:7]:", str4[0:7])
fmt.Println("string(bytes2[0:7]):", string(bytes2[0:7]))
fmt.Println("string(runes2[0:3]):", string(runes2[0:3])) // Correct for Chinese
```

**Notes:**
- `len(str)` returns the number of bytes, not characters.
- Chinese characters are 3 bytes (UTF-8), so use rune slices for correct slicing.

---

## Summary

- Go offers a variety of data types: integers, floats, complex numbers, strings, each with unique features.
- Distinguish between bytes and characters when handling strings and Unicode.
- Type conversion is important for comparisons and calculations.
- Choose between raw and escaped strings as needed.