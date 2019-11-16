package integrit

import "github.com/google/uuid"

type validation interface {
	match(appid uuid.UUID, apiUpload string, agentUpload string) (bool error)
}

/*
RSA((Hash(elf) + Header(elf)), KEYpub) = app_
File = Pack(app_id + app_ + NoHeader(elf))


appid, app_, no_header_elf = UnPack(File)
KEYpri = HTTP(GET /key/pri?appid=appid)
hash, elf_header = RSA(app_, KEYpri)
elf = FillHeader(elf_header, no_header_elf)

if hash = Hash(elf):
    return DONE
else:
    return ERROR
 */