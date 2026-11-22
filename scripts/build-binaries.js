const { execSync } = require('child_process');
const path = require('path');
const fs = require('fs');

const rootDir = path.resolve(__dirname, '..');
const outDir = path.resolve(rootDir, 'nodejs-wrapper', 'bin');

if (!fs.existsSync(outDir)) {
  fs.mkdirSync(outDir, { recursive: true });
}

const targets = [
  { os: 'windows', arch: 'amd64', bin: 'frr-cli-windows-amd64.exe' },
  { os: 'windows', arch: 'arm64', bin: 'frr-cli-windows-arm64.exe' },
  { os: 'linux', arch: 'amd64', bin: 'frr-cli-linux-amd64' },
  { os: 'linux', arch: 'arm64', bin: 'frr-cli-linux-arm64' },
  { os: 'darwin', arch: 'amd64', bin: 'frr-cli-darwin-amd64' },
  { os: 'darwin', arch: 'arm64', bin: 'frr-cli-darwin-arm64' },
];

console.log('🚀 Compilando binarios Go optimizados (-s -w)...');

for (const target of targets) {
  const destPath = path.join(outDir, target.bin);
  console.log(`📦 [${target.os}/${target.arch}] -> ${target.bin}`);
  
  const env = {
    ...process.env,
    GOOS: target.os,
    GOARCH: target.arch,
    CGO_ENABLED: '0'
  };

  try {
    execSync(`go build -ldflags="-s -w" -o "${destPath}" ./cmd/frr-cli`, {
      cwd: rootDir,
      env,
      stdio: 'inherit'
    });
  } catch (err) {
    console.error(`❌ Error al compilar para ${target.os}/${target.arch}:`, err.message);
    process.exit(1);
  }
}

// También crear frr-cli.exe y frr-cli para desarrollo local directo
const localIsWin = process.platform === 'win32';
const localDefaultBin = localIsWin ? 'frr-cli.exe' : 'frr-cli';
const srcTarget = localIsWin ? 'frr-cli-windows-amd64.exe' : 'frr-cli-linux-amd64';
fs.copyFileSync(path.join(outDir, srcTarget), path.join(outDir, localDefaultBin));

console.log('✅ Todos los binarios fueron compilados exitosamente.');
