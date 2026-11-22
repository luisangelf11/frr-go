import * as fs from 'fs';
import * as path from 'path';

export function resolveBinaryPath(): string {
  if (process.env.REPORTERIA_BIN && fs.existsSync(process.env.REPORTERIA_BIN)) {
    return process.env.REPORTERIA_BIN;
  }

  const isWindows = process.platform === 'win32';
  const binName = isWindows ? 'reporteria-cli.exe' : 'reporteria-cli';

  // Buscar en varias rutas relativas estándar
  const possiblePaths = [
    path.resolve(__dirname, '..', '..', 'bin', binName),
    path.resolve(__dirname, '..', 'bin', binName),
    path.resolve(process.cwd(), 'bin', binName),
    path.resolve(process.cwd(), '..', 'bin', binName),
  ];

  for (const p of possiblePaths) {
    if (fs.existsSync(p)) {
      return p;
    }
  }

  // Fallback: nombre directo para confiar en PATH
  return binName;
}
