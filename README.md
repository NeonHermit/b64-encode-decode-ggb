## Search and encode - decode b64 strings from GeoGebra files saved as html  

CLI tool for searching, encoding, decoding and replacing the b64 strings in .html files from GeoGebra.  

- My typical structure: 
  ```
  |_dir_name
    |_index.html
    |_deployggb.js
  ```

- Flags:  
  ```
  --input
  --output
  --zip
  --unzip
  --encode
  --replace
  ```

- Usage:  
  ```
  b64-geogb --input '.\path\to\sourceDir' --output '.\path\to\outDir' 
  ```
  
  Since .ggb is just a .zip you can pass the --unzip flag. Useful if you have to translate large numbers of exercises.
  ```
  b64-geogb --input '.\path\to\sourceDir' --output '.\path\to\outDir' --unzip
  ```
  
  If you unzipped the files, you have to zipp them back with --zip
  ```
  b64-geogb --input '.\path\to\outDir' --output '.\path\to\zippDir' --zip
  ```

  And you can replace the newly edited files back in the .html with --encode & --replace
  ```
  b64-geogb --encode --input '.\path\to\zippDir' --replace '.\path\to\sourceDir'
  ```
