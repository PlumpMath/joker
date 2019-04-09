// This file is generated by generate-std.joke script. Do not edit manually!

package crypto

import (
	. "github.com/candid82/joker/core"
)

var cryptoNamespace = GLOBAL_ENV.EnsureNamespace(MakeSymbol("joker.crypto"))

var hmac_ Proc = func(_args []Object) Object {
	_c := len(_args)
	switch {
	case _c == 3:
		algorithm := ExtractKeyword(_args, 0)
		message := ExtractString(_args, 1)
		key := ExtractString(_args, 2)
		_res := hmacSum(algorithm, message, key)
		return MakeString(_res)

	default:
		PanicArity(_c)
	}
	return NIL
}

func init() {

	cryptoNamespace.ResetMeta(MakeMeta(nil, "Implements common cryptographic and hash functions.", "1.0"))

	cryptoNamespace.InternVar("hmac", hmac_,
		MakeMeta(
			NewListFrom(NewVectorFrom(MakeSymbol("algorithm"), MakeSymbol("message"), MakeSymbol("key"))),
			`Returns HMAC signature for message and key using specified algorithm.
  Algorithm is one of the following: :sha1, :sha224, :sha256, :sha384, :sha512.`, "1.0"))

}