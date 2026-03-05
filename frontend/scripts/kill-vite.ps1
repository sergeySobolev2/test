Get-Process -Name node -ErrorAction SilentlyContinue |
  Where-Object { $_.Path -like '*node_modules*vite*' -or $_.Path -like '*vite*' -or $_.ProcessName -eq 'node' } |
  ForEach-Object {
    try { Stop-Process -Id $_.Id -Force -ErrorAction Stop } catch {}
  }
