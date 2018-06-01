all:
	go build -o main draw.go image.go main.go matrix.go vector.go stack.go lighting.go lexer.go token.go parser.go command.go

clean:
	rm *.ppm
	rm *.png
	rm main
