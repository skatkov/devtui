#!/bin/bash

# Auto-generated script to run all tape demos
# Make sure you have vhs installed: go install github.com/charmbracelet/vhs@latest

set -e

# Change to docs directory if not already there
cd "$(dirname "$0")"

echo "Generating screenshots for all TUI modules..."

echo "Generating main menu demo..."
vhs tapes/demo-main.tape
sleep 1

echo "Generating screenshot for Base64 Decoder..."
vhs tapes/demo-base64-decoder.tape
sleep 1

echo "Generating screenshot for Base64 Encoder..."
vhs tapes/demo-base64-encoder.tape
sleep 1

echo "Generating screenshot for Cron..."
vhs tapes/demo-cron.tape
sleep 1

echo "Generating screenshot for Css..."
vhs tapes/demo-css.tape
sleep 1

echo "Generating screenshot for Csv2md..."
vhs tapes/demo-csv2md.tape
sleep 1

echo "Generating screenshot for Csvjson..."
vhs tapes/demo-csvjson.tape
sleep 1

echo "Generating screenshot for Graphql Query..."
vhs tapes/demo-graphql-query.tape
sleep 1

echo "Generating screenshot for Html..."
vhs tapes/demo-html.tape
sleep 1

echo "Generating screenshot for Json..."
vhs tapes/demo-json.tape
sleep 1

echo "Generating screenshot for Jsonstruct..."
vhs tapes/demo-jsonstruct.tape
sleep 1

echo "Generating screenshot for Jsontoml..."
vhs tapes/demo-jsontoml.tape
sleep 1

echo "Generating screenshot for Markdown..."
vhs tapes/demo-markdown.tape
sleep 1

echo "Generating screenshot for Numbers..."
vhs tapes/demo-numbers.tape
sleep 1

echo "Generating screenshot for Toml..."
vhs tapes/demo-toml.tape
sleep 1

echo "Generating screenshot for Tomljson..."
vhs tapes/demo-tomljson.tape
sleep 1

echo "Generating screenshot for Tsv2md..."
vhs tapes/demo-tsv2md.tape
sleep 1

echo "Generating screenshot for Uuiddecode..."
vhs tapes/demo-uuiddecode.tape
sleep 1

echo "Generating screenshot for Uuidgenerate..."
vhs tapes/demo-uuidgenerate.tape
sleep 1

echo "Generating screenshot for Xml..."
vhs tapes/demo-xml.tape
sleep 1

echo "Generating screenshot for Yaml..."
vhs tapes/demo-yaml.tape
sleep 1

echo "Generating screenshot for Yamlstruct..."
vhs tapes/demo-yamlstruct.tape
sleep 1

echo "All screenshots generated successfully!"
