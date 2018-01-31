// This file is generated by generate-std.joke script. Do not edit manually!

package os

import (
  "os"
  . "github.com/candid82/joker/core"
)

var osNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.os"))

var args_ Proc = func(args []Object) Object {
  c := len(args)
  switch  {
  case c == 0:
    
    
    res := commandArgs()
    return res

  default:
    PanicArity(c)
  }
  return NIL
}

var env_ Proc = func(args []Object) Object {
  c := len(args)
  switch  {
  case c == 0:
    
    
    res := env()
    return res

  default:
    PanicArity(c)
  }
  return NIL
}

var exit_ Proc = func(args []Object) Object {
  c := len(args)
  switch  {
  case c == 1:
    
    code := ExtractInt(args, 0)
    res := NIL; os.Exit(code)
    return res

  default:
    PanicArity(c)
  }
  return NIL
}

var ls_ Proc = func(args []Object) Object {
  c := len(args)
  switch  {
  case c == 1:
    
    dirname := ExtractString(args, 0)
    res := readDir(dirname)
    return res

  default:
    PanicArity(c)
  }
  return NIL
}

var mkdir_ Proc = func(args []Object) Object {
  c := len(args)
  switch  {
  case c == 2:
    
    name := ExtractString(args, 0)
    perm := ExtractInt(args, 1)
    res := mkdir(name, perm)
    return res

  default:
    PanicArity(c)
  }
  return NIL
}

var sh_ Proc = func(args []Object) Object {
  c := len(args)
  switch  {
  case true:
    CheckArity(args, 1,999)
    name := ExtractString(args, 0)
    arguments := ExtractStrings(args, 1)
    res := sh(name, arguments)
    return res

  default:
    PanicArity(c)
  }
  return NIL
}


func init() {

osNamespace.ResetMeta(MakeMeta(nil, "Provides a platform-independent interface to operating system functionality.", "1.0"))

osNamespace.InternVar("args", args_,
  MakeMeta(
    NewListFrom(NewVectorFrom()),
    `Returns a sequence of the command line arguments, starting with the program name (normally, joker).`, "1.0"))

osNamespace.InternVar("env", env_,
  MakeMeta(
    NewListFrom(NewVectorFrom()),
    `Returns a map representing the environment.`, "1.0"))

osNamespace.InternVar("exit", exit_,
  MakeMeta(
    NewListFrom(NewVectorFrom(MakeSymbol("code"))),
    `Causes the current program to exit with the given status code.`, "1.0"))

osNamespace.InternVar("ls", ls_,
  MakeMeta(
    NewListFrom(NewVectorFrom(MakeSymbol("dirname"))),
    `Reads the directory named by dirname and returns a list of directory entries sorted by filename.`, "1.0"))

osNamespace.InternVar("mkdir", mkdir_,
  MakeMeta(
    NewListFrom(NewVectorFrom(MakeSymbol("name"), MakeSymbol("perm"))),
    `Creates a new directory with the specified name and permission bits.`, "1.0"))

osNamespace.InternVar("sh", sh_,
  MakeMeta(
    NewListFrom(NewVectorFrom(MakeSymbol("name"), MakeSymbol("&"), MakeSymbol("arguments"))),
    `Executes the named program with the given arguments. Returns a map with the following keys:
      :success - whether or not the execution was successful,
      :out - string capturing stdout of the program,
      :err - string capturing stderr of the program.`, "1.0"))

}