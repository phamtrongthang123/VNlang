package evaluator

import (
	"math/big"
	"strings"
	"testing"
	"vnlang/lexer"
	"vnlang/object"
	"vnlang/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 + -50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestEvalFloatExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"3.14", 3.14},
		{"-3.14", -3.14},
		{"0.5 + .5 + 1.5 + 2.5 - 10.0", -5},
		{"2.0 * 2.2 * 2.4 * 2.6 * 2.8", 76.8768},
		{"-50.0 + 100.0 + -50.0", 0},
		{"5.5 * 2.1 + 10.1", 21.65},
		{"5.3 + 2.1 * 10.3", 26.93},
		{"20.7 + 2.4 * -10.3", -4.02},
		{"50.2 / 2.1 * 2.1 + 10.8", 61},
		{"2.2 * (5.5 + 10.2)", 34.54},
		{"3.3 * 3.3 * 3.3 + 10.3", 46.237},
		{"3.2 * (3.3 * 3.4) + 10.6", 46.504},
		{"(5.1 + 10.2 * 2.9 + 15.8 / 3.7) * 2.6 + -10.5", 90.7707027027},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testFloatObject(t, evaluated, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"đúng", true},
		{"sai", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"đúng == đúng", true},
		{"sai == sai", true},
		{"đúng == sai", false},
		{"đúng != sai", true},
		{"sai != đúng", true},
		{"(1 < 2) == đúng", true},
		{"(1 < 2) == sai", false},
		{"(1 > 2) == đúng", false},
		{"(1 > 2) == sai", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestBangOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!đúng", false},
		{"!sai", true},
		{"!5", false},
		{"!!đúng", true},
		{"!!sai", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"nếu (đúng) { 10 }", 10},
		{"nếu (sai) { 10 }", nil},
		{"nếu (1) { 10 }", 10},
		{"nếu (1 < 2) { 10 }", 10},
		{"nếu (1 > 2) { 10 }", nil},
		{"nếu (1 > 2) { 10 } ngược_lại { 20 }", 20},
		{"nếu (1 < 2) { 10 } ngược_lại { 20 }", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"trả_về 10;", 10},
		{"trả_về 10; 9;", 10},
		{"trả_về 2 * 5; 9;", 10},
		{"9; trả_về 2 * 5; 9;", 10},
		{"nếu (10 > 1) { trả_về 10; }", 10},
		{
			`
nếu (10 > 1) {
  nếu (10 > 1) {
    trả_về 10;
  }

  trả_về 1;
}
`,
			10,
		},
		{
			`
đặt f = hàm(x) {
  trả_về x;
  x + 10;
};
f(10);`,
			10,
		},
		{
			`
đặt f = hàm(x) {
   đặt kết_quả = x + 10;
   trả_về kết_quả;
   trả_về 10;
};
f(10);`,
			20,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + đúng;",
			"kiểu không tương thích: SỐ_NGUYÊN + BOOLEAN",
		},
		{
			"5 + đúng; 5;",
			"kiểu không tương thích: SỐ_NGUYÊN + BOOLEAN",
		},
		{
			"-đúng",
			"toán tử lạ: -BOOLEAN",
		},
		{
			"đúng + sai;",
			"toán tử lạ: BOOLEAN + BOOLEAN",
		},
		{
			"đúng + sai + đúng + sai;",
			"toán tử lạ: BOOLEAN + BOOLEAN",
		},
		{
			"5; đúng + sai; 5",
			"toán tử lạ: BOOLEAN + BOOLEAN",
		},
		{
			`"Hello" - "World"`,
			"toán tử lạ: XÂU - XÂU",
		},
		{
			"nếu (10 > 1) { đúng + sai; }",
			"toán tử lạ: BOOLEAN + BOOLEAN",
		},
		{
			`
nếu (10 > 1) {
  nếu (10 > 1) {
    trả_về đúng + sai;
  }

  trả_về 1;
}
`,
			"toán tử lạ: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"không tìm thấy tên định danh: foobar",
		},
		{
			`{"name": "Monkey"}[hàm(x) { x }];`,
			"không thể dùng như khóa băm: HÀM",
		},
		{
			`999[1]`,
			"toán tử chỉ mục không hỗ trợ cho: SỐ_NGUYÊN",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"đặt a = 5; a;", 5},
		{"đặt a = 5 * 5; a;", 25},
		{"đặt a = 5; đặt b = a; b;", 5},
		{"đặt a = 5; đặt b = a; đặt c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "hàm(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v",
			fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	expectedBody := "(x + 2) "

	if fn.Body.String() != expectedBody {
		t.Fatalf("body is not %q. got=%q", expectedBody, fn.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"đặt identity = hàm(x) { x; }; identity(5);", 5},
		{"đặt identity = hàm(x) { trả_về x; }; identity(5);", 5},
		{"đặt double = hàm(x) { x * 2; }; double(5);", 10},
		{"đặt add = hàm(x, y) { x + y; }; add(5, 5);", 10},
		{"đặt add = hàm(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"hàm(x) { x; }(5)", 5},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestEnclosingEnvironments(t *testing.T) {
	input := `
đặt first = 10;
đặt second = 10;
đặt third = 10;

đặt ourFunction = hàm(first) {
  đặt second = 20;

  first + second + third;
};

ourFunction(20) + first + second;`

	testIntegerObject(t, testEval(input), 70)
}

func TestClosures(t *testing.T) {
	input := `
đặt newAdder = hàm(x) {
  hàm(y) { x + y };
};

đặt addTwo = newAdder(2);
addTwo(2);`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatenation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{`độ_dài("")`, 0},
		{`độ_dài("four")`, 4},
		{`độ_dài("hello world")`, 11},
		{`độ_dài(1)`, "Tham số truyền vào `độ_dài` không được hỗ trợ lấy độ dài (chỉ có Mảng hoặc Chuỗi được hỗ trợ), kiểu tham số SỐ_NGUYÊN."},
		{`độ_dài("one", "two")`, "Sai số lượng tham số truyền vào. nhận được = 2, mong muốn = 1"},
		{`độ_dài([1, 2, 3])`, 3},
		{`độ_dài([])`, 0},
		{`in_ra("hello", "world!")`, nil},
		{`đầu([1, 2, 3])`, 1},
		{`đầu([])`, nil},
		{`đầu(1)`, "Tham số truyền vào hàm lấy `đầu` của mảng phải thuộc kiểu Mảng. Nhận được kiểu SỐ_NGUYÊN"},
		{`đuôi([1, 2, 3])`, 3},
		{`đuôi([])`, nil},
		{`đuôi(1)`, "Tham số truyền vào hàm lấy `đuôi` của mảng phải thuộc kiểu Mảng. Nhận được kiểu SỐ_NGUYÊN"},
		{`trừ_đầu([1, 2, 3])`, []int{2, 3}},
		{`trừ_đầu([])`, nil},
		{`đẩy([], 1)`, []int{1}},
		{`đẩy(1, 1)`, "Tham số truyền vào hàm lấy `đẩy` của mảng phải thuộc kiểu Mảng. Nhận được kiểu SỐ_NGUYÊN"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		switch expected := tt.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case nil:
			testNullObject(t, evaluated)
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)",
					evaluated, evaluated)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q",
					expected, errObj.Message)
			}
		case []int:
			array, ok := evaluated.(*object.Array)
			if !ok {
				t.Errorf("obj not Array. got=%T (%+v)", evaluated, evaluated)
				continue
			}

			if len(array.Elements) != len(expected) {
				t.Errorf("wrong num of elements. want=%d, got=%d",
					len(expected), len(array.Elements))
				continue
			}

			for i, expectedElem := range expected {
				testIntegerObject(t, array.Elements[i], int64(expectedElem))
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d",
			len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func TestArrayIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			"[1, 2, 3][0]",
			1,
		},
		{
			"[1, 2, 3][1]",
			2,
		},
		{
			"[1, 2, 3][2]",
			3,
		},
		{
			"đặt i = 0; [1][i];",
			1,
		},
		{
			"[1, 2, 3][1 + 1];",
			3,
		},
		{
			"đặt myArray = [1, 2, 3]; myArray[2];",
			3,
		},
		{
			"đặt myArray = [1, 2, 3]; myArray[0] + myArray[1] + myArray[2];",
			6,
		},
		{
			"đặt myArray = [1, 2, 3]; đặt i = myArray[0]; myArray[i]",
			2,
		},
		{
			"[1, 2, 3][3]",
			nil,
		},
		{
			"[1, 2, 3][-1]",
			nil,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestHashLiterals(t *testing.T) {
	input := `đặt two = "two";
	{
		"one": 10 - 9,
		two: 1 + 1,
		"thr" + "ee": 6 / 2,
		4: 4,
		đúng: 5,
		sai: 6
	}`

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Hash)
	if !ok {
		t.Fatalf("Eval didn't return Hash. got=%T (%+v)", evaluated, evaluated)
	}

	expected := map[object.HashKey]int64{
		(&object.String{Value: "one"}).HashKey():          1,
		(&object.String{Value: "two"}).HashKey():          2,
		(&object.String{Value: "three"}).HashKey():        3,
		(&object.Integer{Value: big.NewInt(4)}).HashKey(): 4,
		TRUE.HashKey():  5,
		FALSE.HashKey(): 6,
	}

	if len(result.Pairs) != len(expected) {
		t.Fatalf("Hash has wrong num of pairs. got=%d", len(result.Pairs))
	}

	for expectedKey, expectedValue := range expected {
		pair, ok := result.Pairs[expectedKey]
		if !ok {
			t.Errorf("no pair for given key in Pairs")
		}

		testIntegerObject(t, pair.Value, expectedValue)
	}
}

func TestHashIndexExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			`{"foo": 5}["foo"]`,
			5,
		},
		{
			`{"foo": 5}["bar"]`,
			nil,
		},
		{
			`đặt key = "foo"; {"foo": 5}[key]`,
			5,
		},
		{
			`{}["foo"]`,
			nil,
		},
		{
			`{5: 5}[5]`,
			5,
		},
		{
			`{đúng: 5}[đúng]`,
			5,
		},
		{
			`{sai: 5}[sai]`,
			5,
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func testEval(input string) object.Object {
	l := lexer.New(strings.NewReader(input), "<test>")
	p := parser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(object.NewCallStack(), program, env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if !result.Value.IsInt64() || result.Value.Int64() != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}

	return true
}

func testFloatObject(t *testing.T, obj object.Object, expected float64) bool {
	result, ok := obj.(*object.Float)
	if !ok {
		t.Errorf("object is not Float. got=%T (%+v)", obj, obj)
		return false
	}
	if !(result.Value-expected <= 1e-9) {
		t.Errorf("object has wrong value. got=%f, want=%f",
			result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
