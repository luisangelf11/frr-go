import { spawn, ChildProcess } from 'child_process';
import { resolveBinaryPath } from './binary-resolver';
import { ReportSpec, ExportFormat } from './types';

interface PendingRequest {
  resolve: (data: Buffer) => void;
  reject: (err: Error) => void;
}

export class GoWorkerManager {
  private static instance: GoWorkerManager | null = null;
  private process: ChildProcess | null = null;
  private buffer: Buffer = Buffer.alloc(0);
  private queue: PendingRequest[] = [];
  private binaryPath: string;

  constructor(customBinPath?: string) {
    this.binaryPath = customBinPath || resolveBinaryPath();
    this.initProcess();
  }

  public static getInstance(customBinPath?: string): GoWorkerManager {
    if (!GoWorkerManager.instance) {
      GoWorkerManager.instance = new GoWorkerManager(customBinPath);
    }
    return GoWorkerManager.instance;
  }

  private initProcess() {
    if (this.process) {
      return;
    }

    this.process = spawn(this.binaryPath, ['--worker'], {
      stdio: ['pipe', 'pipe', 'inherit'],
    });

    this.process.stdout?.on('data', (chunk: Buffer) => {
      this.buffer = Buffer.concat([this.buffer, chunk]);
      this.drainBuffer();
    });

    this.process.on('error', (err: Error) => {
      while (this.queue.length > 0) {
        const req = this.queue.shift();
        req?.reject(new Error(`Error en el proceso worker Go: ${err.message}`));
      }
      this.process = null;
    });

    this.process.on('close', (code: number | null) => {
      while (this.queue.length > 0) {
        const req = this.queue.shift();
        req?.reject(new Error(`El worker de Go terminó con código ${code}`));
      }
      this.process = null;
    });

    // Limpieza limpia al cerrar la app Node.js
    const cleanup = () => this.destroy();
    if (typeof process !== 'undefined') {
      process.once('exit', cleanup);
      process.once('SIGINT', cleanup);
      process.once('SIGTERM', cleanup);
    }
  }

  private drainBuffer() {
    while (this.buffer.length >= 5) {
      const dataLen = this.buffer.readUInt32BE(0);
      const totalMsgLen = 4 + 1 + dataLen;

      if (this.buffer.length < totalMsgLen) {
        break; // Esperar más chunks de datos
      }

      const status = this.buffer.readUInt8(4);
      const payloadData = this.buffer.subarray(5, totalMsgLen);
      this.buffer = this.buffer.subarray(totalMsgLen);

      const pending = this.queue.shift();
      if (pending) {
        if (status === 0) {
          pending.resolve(payloadData);
        } else {
          pending.reject(new Error(payloadData.toString('utf-8')));
        }
      }
    }
  }

  public execute(spec: ReportSpec, format: ExportFormat, delimiter?: string): Promise<Buffer> {
    return new Promise((resolve, reject) => {
      if (!this.process || this.process.killed) {
        this.initProcess();
      }

      if (!this.process || !this.process.stdin) {
        return reject(new Error('No se pudo inicializar el proceso worker de Go'));
      }

      const payload = JSON.stringify({
        format,
        delimiter: delimiter || ',',
        spec,
      });

      const payloadBuf = Buffer.from(payload, 'utf-8');
      const headerBuf = Buffer.alloc(4);
      headerBuf.writeUInt32BE(payloadBuf.length, 0);

      this.queue.push({ resolve, reject });

      this.process.stdin.write(headerBuf);
      this.process.stdin.write(payloadBuf);
    });
  }

  public destroy() {
    if (this.process) {
      try {
        this.process.kill();
      } catch (_) {}
      this.process = null;
    }
  }
}
