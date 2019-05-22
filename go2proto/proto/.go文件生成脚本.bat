md "./go/"
for /f "delims=" %%A in ('dir /b *.proto') do  protoc --go_out=plugins=grpc:./go/  ./%%A
pause