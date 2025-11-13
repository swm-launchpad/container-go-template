import { Controller, Get, Post, Delete, Param, Body, HttpCode, HttpStatus, NotFoundException, BadRequestException } from '@nestjs/common';

interface Item {
  id: number;
  name: string;
  description: string;
  createdAt: string;
}

interface CreateItemDto {
  name: string;
  description?: string;
}

@Controller()
export class AppController {
  private items: Item[] = [
    { id: 1, name: 'Sample Item 1', description: 'First sample item', createdAt: new Date().toISOString() },
    { id: 2, name: 'Sample Item 2', description: 'Second sample item', createdAt: new Date().toISOString() }
  ];
  private nextId = 3;
  private requestCount = 0;
  private readonly startTime = Date.now();

  constructor() {
    // Request counter could be implemented via middleware/interceptor in production
  }

  @Get()
  getRoot() {
    this.requestCount++;
    return {
      message: 'NestJS Test API',
      version: '1.0.0',
      requestCount: this.requestCount,
      timestamp: new Date().toISOString(),
      endpoints: [
        'GET /',
        'GET /health',
        'GET /items',
        'POST /items',
        'GET /items/:id',
        'DELETE /items/:id'
      ]
    };
  }

  @Get('health')
  getHealth() {
    this.requestCount++;
    const uptime = Math.floor((Date.now() - this.startTime) / 1000);
    return {
      status: 'healthy',
      uptime: `${uptime} seconds`,
      timestamp: new Date().toISOString()
    };
  }

  @Get('items')
  getItems() {
    this.requestCount++;
    return {
      success: true,
      count: this.items.length,
      data: this.items
    };
  }

  @Post('items')
  @HttpCode(HttpStatus.CREATED)
  createItem(@Body() createItemDto: CreateItemDto) {
    this.requestCount++;

    if (!createItemDto.name) {
      throw new BadRequestException('Name is required');
    }

    const newItem: Item = {
      id: this.nextId++,
      name: createItemDto.name,
      description: createItemDto.description || '',
      createdAt: new Date().toISOString()
    };

    this.items.push(newItem);

    return {
      success: true,
      data: newItem
    };
  }

  @Get('items/:id')
  getItem(@Param('id') id: string) {
    this.requestCount++;
    const itemId = parseInt(id, 10);
    const item = this.items.find(i => i.id === itemId);

    if (!item) {
      throw new NotFoundException('Item not found');
    }

    return {
      success: true,
      data: item
    };
  }

  @Delete('items/:id')
  deleteItem(@Param('id') id: string) {
    this.requestCount++;
    const itemId = parseInt(id, 10);
    const index = this.items.findIndex(i => i.id === itemId);

    if (index === -1) {
      throw new NotFoundException('Item not found');
    }

    const deletedItem = this.items.splice(index, 1)[0];

    return {
      success: true,
      message: 'Item deleted successfully',
      data: deletedItem
    };
  }
}
