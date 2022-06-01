go build -buildmode=c-archive main.go clibh.go

gcc degate.def main.a -shared -lwinmm -lWs2_32 -o degate.dll -Wl,--out-implib,degate.lib