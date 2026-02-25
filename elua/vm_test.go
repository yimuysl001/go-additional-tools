package elua

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/yuin/gopher-lua"
	"layeh.com/gopher-luar"
	"testing"
)

type Song struct {
	Title  string
	Artist string
}

func (s Song) Sum(a, b string) string {
	return a + b
}

func TestExecScript(t *testing.T) {

	//RegisterFunc("sum", func(a, b string) string {
	//	return a + b
	//})
	RegisterTypeFunc("Song", Song{})

	v, err := ExecScript(gctx.GetInitCtx(), `

s = Song()
b = s:Sum( "5" , params.u)
s.Title= "测试"
s.Artist= b
return s
`, map[string]any{
		"u": "1",
	})
	fmt.Println(v, err)
}

func TestLua1(t *testing.T) {
	const code = `
	print(sum(1, 2, 3, 4, 5))
	`

	L := lua.NewState()
	defer L.Close()

	sum := func(L *luar.LState) int {
		total := 0
		for i := 1; i <= L.GetTop(); i++ {
			total += L.CheckInt(i)
		}
		L.Push(lua.LNumber(total))
		return 1
	}

	L.SetGlobal("sum", luar.New(L, sum))

	if err := L.DoString(code); err != nil {
		panic(err)
	}
}

func TestLua2(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	u := &User{
		Name: "Tim",
	}
	L.SetGlobal("u", luar.New(L, u))
	if err := L.DoString(`script`); err != nil {
		panic(err)
	}

	fmt.Println("Lua set your token to:", u.Token())
	// Output:
	// Hello from Lua, Tim!
	// Lua set your token to: 12345

}

func TestLua3(t *testing.T) {

	L := lua.NewState()
	defer L.Close()

	u := &User{
		Name: "Tim",
	}
	L.SetGlobal("u", luar.New(L, u))
	var m = make(map[string]any)
	m = map[string]any{
		"CtxId": gctx.CtxId,
	}
	L.SetGlobal("u", luar.New(L, u))
	L.SetGlobal("m", luar.New(L, m))
	L.SetGlobal("log", luar.New(L, g.Log()))
	L.SetGlobal("ctx", luar.New(L, gctx.New()))
	if err := L.DoString(`
log:Info(ctx,"Hello from Lua, " .. u.Name .. "!")
log:Info(ctx,"Hello from Lua, " .. m.CtxId(ctx) .. "!")
print("Hello from Lua, " .. u.Name .. "!")
u:SetToken("12345")
m.a="1"

`); err != nil {
		panic(err)
	}

	fmt.Println("Lua set your token to:", u.Token())
	// Output:
	// Hello from Lua, Tim!
	// Lua set your token to: 12345
	fmt.Println("Lua m your token to:", m)

}

func TestFor(t *testing.T) {
	const (
		cript = `-- 自定义计算表中最大键值函数 table_maxn，即返回表最大键值
function table_maxn(t)
    local mn = 0
    for k, _ in pairs(t) do
        if type(k) == "number" and k > mn then
            mn = k
        end
    end
    return mn
end

-- 两表相加操作
mytable = setmetatable({ 1, 2, 3 }, {
  __add = function(mytable, newtable)
    local max_key_mytable = table_maxn(mytable)
    for i = 1, table_maxn(newtable) do
      table.insert(mytable, max_key_mytable + i, newtable[i])
    end
    return mytable
  end
})

secondtable = {4, 5, 6}

mytable = mytable + secondtable

for k, v in ipairs(mytable) do
    print(k, v)
end`
	)

	L := lua.NewState()
	defer L.Close()

	err := L.DoString(cript)

	fmt.Println(err)

}

func TestCoroutine(t *testing.T) {
	const (
		cript = `
import("module as m")
 
print(m.constant)

function foo (a)
    print("foo 函数输出", a)
    return coroutine.yield(2 * a) -- 返回  2*a 的值
end
 
co = coroutine.create(function (a , b)
    print("第一次协同程序执行输出", a, b) -- co-body 1 10
    local r = foo(a + 1)
     
    print("第二次协同程序执行输出", r)
    local r, s = coroutine.yield(a + b, a - b)  -- a，b的值为第一次调用协同程序时传入
     
    print("第三次协同程序执行输出", r, s)
    return b, "结束协同程序"                   -- b的值为第二次调用协同程序时传入
end)
       
print("main", coroutine.resume(co, 1, 10)) -- true, 4
print("--分割线----")
print("main", coroutine.resume(co, "r")) -- true 11 -9
print("---分割线---")
print("main", coroutine.resume(co, "x", "y")) -- true 10 end
print("---分割线---")
print("main", coroutine.resume(co, "x", "y")) -- cannot resume dead coroutine
print("---分割线---")
`
	)

	L := lua.NewState()
	defer L.Close()

	RegisterFuncString("module", `module = {}
 
-- 定义一个常量
module.constant = "这是一个常量"
 
-- 定义一个函数
function module.func1()
    io.write("这是一个公有函数！\n")
end
 
local function func2()
    print("这是一个私有函数！")
end
  
function module.func3()
    func2()
end
 
return module`)

	L.SetGlobal("import", luar.New(L, Require))

	err := L.DoString(cript)

	fmt.Println(err)

}

func ExampleNewType() {
	L := lua.NewState()
	defer L.Close()

	type Song struct {
		Title  string
		Artist string
	}

	L.SetGlobal("Song", luar.NewType(L, Song{}))
	if err := L.DoString(`
		s = Song()
		s.Title = "Montana"
		s.Artist = "Tycho"
		print(s.Artist .. " - " .. s.Title)
	`); err != nil {
		panic(err)
	}
	// Output:
	// Tycho - Montana
}

type User struct {
	Name  string
	token string
}

func (u *User) SetToken(t string) {
	u.token = t
}

func (u *User) Token() string {
	return u.token
}

const script = `
print("Hello from Lua, " .. u.Name .. "!")
u:SetToken("12345")
`
