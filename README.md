Here's a simple README for the Go program:

---

# Golang Paprika Recipe Archive to Json 

This program extracts specified fields from gzipped JSON files in a directory (which is the paprikarecipies format) or a zip archive and writes the extracted fields to standard output.

## Dependencies

This program has no external dependencies, as it only uses packages from the Go standard library.

## Installation

1. Install [Go](https://golang.org/dl/) if you haven't already.
2. Clone or download this repository to your local machine.
3. Navigate to the directory containing the program using a terminal or command prompt.

## Usage

To build and run the program, use the following command:

```sh
go run main.go -input /path/to/your/input/directory/or/zip -fields field1,field2,field3
```

Replace `/path/to/your/input/directory/or/zip` with the path to your input directory containing gzipped JSON files or a zip archive containing the gzipped JSON files. Replace `field1,field2,field3` with a comma-separated list of field names you want to extract from the JSON files.

### Command-Line Flags

- `-input`: The path to the input directory containing gzipped JSON files or a zip archive containing gzipped JSON files (default: `./input`).
- `-fields`: A comma-separated list of field names to extract from the JSON files (default: `name,categories`).

### Example

Assuming you have a directory named "recipes" containing gzipped JSON files and you want to extract the "name" and "categories" fields from each file, you can run the following command:

```sh
go run main.go -input recipes -fields name,categories
```

The program will process the gzipped JSON files, extract the specified fields, and write the output JSON to standard output.
