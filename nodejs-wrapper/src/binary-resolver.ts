import * as fs from 'fs';
import * as path from 'path';

export function resolveBinaryPath(): string {
  // 1. Variable de entorno explícita
  const envBin = process.env.FRR_BIN || process.env.REPORTERIA_BIN;
  if (envBin && fs.existsSync(envBin)) {
    return envBin;
  }

  const isWindows = process.platform === 'win32';
  const ext = isWindows ? '.exe' : '';

  // Determinar arquitectura de Go
  let goArch = 'amd64';
  if (process.arch === 'arm64') {
    goArch = 'arm64';
  } else if (process.arch === 'x64') {
    goArch = 'amd64';
  } else if (process.arch === 'ia32') {
    goArch = '386';
  }

  // Determinar OS de Go
  let goOs = 'linux';
  if (process.platform === 'win32') {
    goOs = 'windows';
  } else if (process.platform === 'darwin') {
    goOs = 'darwin';
  }

  const platformSpecificBin = `frr-cli-${goOs}-${goArch}${ext}`;
  const genericBin = isWindows ? 'frr-cli.exe' : 'frr-cli';

  // Buscar en posibles directorios
  const searchDirs = [
    path.resolve(__dirname, '..', 'bin'),
    path.resolve(__dirname, '..', '..', 'bin'),
    path.resolve(process.cwd(), 'bin'),
    path.resolve(process.cwd(), 'node_modules', 'frr-node', 'bin'),
  ];

  for (const dir of searchDirs) {
    // 1. Intentar el binario específico para la plataforma
    const specificPath = path.join(dir, platformSpecificBin);
    if (fs.existsSync(specificPath)) {
      ensureExecutable(specificPath);
      return specificPath;
    }

    // 2. Intentar el genérico
    const genericPath = path.join(dir, genericBin);
    if (fs.existsSync(genericPath)) {
      ensureExecutable(genericPath);
      return genericPath;
    }
  }

  // Fallback a PATH del sistema
  return genericBin;
}

function ensureExecutable(filePath: string): void {
  if (process.platform !== 'win32') {
    try {
      const stat = fs.statSync(filePath);
      // Si no tiene permisos de ejecución para el dueño (0o100)
      if ((stat.mode & 0o111) === 0) {
        fs.chmodSync(filePath, 0o755);
      }
    } catch {
      // Ignorar si no se pueden modificar permisos
    }
  }
}
