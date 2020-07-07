# VNlang

Là ngôn ngữ lập trình với syntax là tiếng Việt. Được phát triển dựa trên ngôn ngữ Monkey từ quyển **Writing An Interpreter In Go**.

## Description

Các tính năng chính:

- C-like syntax
- variable bindings
- integers and booleans
- arithmetic expressions
- built-in functions
- first-class and higher-order functions
- closures
- a string data structure
- an array data structure
- a hash data structure

## Installation

## Usage

### REPL

```
cd src
go run .\vnlang\main.go
```

Một số mẫu sử dụng:

```
>> đặt a = [1,2,323,4]
>> độ_dài(a)
4
>> đầu(a)
1
>> a[3]
4
>> a[2]
323
```

```
>> đặt chuỗi = "asdasdsad"
>> chuỗi[1]
LỖI: toán tử chỉ mục không hỗ trợ cho: XÂU
>> in_ra(chuỗi)
asdasdsad
null
>> đặt c2 = "test here"
>> chuỗi + " " + c2
asdasdsad test here
```

```
>> đặt fi = hàm(b) { nếu (b==0){trả_về 0;} ngược_lại { nếu (b==1) {trả_về 1;} ngược_lại {trả_về fi(b-1) + fi(b-2); } } }
>> fi
hàm (b) {
nếu (b == 0) trả_về 0; ngược_lại nếu (b == 1) trả_về 1; ngược_lại trả_về (fi((b - 1)) + fi((b - 2)));
}
>> fi(1)
1
>> fi(4)
3
>> fi(9)
34
```

### Run script

#### HTTP example 

Chạy theo hướng dẫn ở repo này: https://github.com/phamtrongthang123/software_design_final_term_project



## License

[MIT](https://choosealicense.com/licenses/mit/)
