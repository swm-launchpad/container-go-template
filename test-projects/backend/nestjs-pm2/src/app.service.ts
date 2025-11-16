import { Injectable } from '@nestjs/common';

@Injectable()
export class AppService {
  getHello() {
    return {
      message: 'Hello from NestJS with Bun!',
      package_manager: 'bun',
      framework: 'NestJS',
    };
  }
}
