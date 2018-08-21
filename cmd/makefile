X_PROD = -X main.prod=1

X_DEBUG = -X main.debug=1

H_WINDOWS_GUI = -H windowsgui

LDFLAGS= -linkmode internal -X 'main.buildtime=%TIME%'

GO_BUILD = go build -i -ldflags

prod:
    ${GO_BUILD} "${LDFLAGS} ${H_WINDOWS_GUI} ${X_PROD}" -o build/amd64/prod.exe

debug:
    ${GO_BUILD}  "${LDFLAGS} ${X_DEBUG}" -o build/amd64/debug.exe


prod32:
    set GOARCH=386 \
    ${GO_BUILD} "${LDFLAGS} ${H_WINDOWS_GUI} ${X_PROD}" -o build/386/prod.exe \
    set GOARCH=amd64

debug32:
    set GOARCH=386 \
    ${GO_BUILD}  "${LDFLAGS} ${X_DEBUG}" -o build/386/debug.exe \
    set GOARCH=amd64