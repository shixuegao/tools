package sxg.main;

import java.nio.ByteBuffer;

public class BitMap {
	public static final long LimitMaxValue = 0xffffffffl;
	long max;
	long offset;
	ByteBuffer buffer;
	
	public BitMap(long max, long offset) throws IllegalArgumentException {
		long value = normalized(max, offset);
		if (value == -1) {
			String errInfo = String.format("非法的值[%d], 值必须在[%d, %d]之间", max, offset+1, offset+LimitMaxValue);
			throw new IllegalArgumentException(errInfo);
		}
		grow(max);
	}
	
	//offset + 1 ~ offset + LimitMaxValue
	private long normalized(long value, long offset) {
		value = value - offset;
		if (value > 0 && value <= LimitMaxValue) {
			return value;
		} else {
			return -1;
		}
	}
	
	private void grow(long max) {
		if (max <= this.max) {
			return;
		}
		
		this.max = max;
		int index = (int)(max >>> 3);
		byte mod = (byte)(max & 0x07);
		if (mod == 0) {
			index -= 1;
		}
		ByteBuffer newBuffer = ByteBuffer.allocate(index + 1);
		if (buffer != null) {
			buffer.position(0);				
			newBuffer.put(buffer);
		}
		buffer = newBuffer;
	}
	
	private byte mask(byte mod) {
		if (mod == 0) {
			return -128;
		} else {
			return (byte)(0x01 << (mod - 1));
		}
	}
	
	public long offset() {
		return offset;
	}
	
	public boolean existed(long value) {
		value = normalized(value, offset);
		if (value < 1 || value > max) {
			return false;
		}
		int index = (int)(value >>> 3);
		byte mod = (byte)(value & 0x07);
		if (mod == 0) {
			index -= 1;
		}
		short mask = mask(mod);
		
		buffer.position(index);
		return (buffer.get() & mask) != 0;
	}
	
	public void put(long value) throws IllegalArgumentException {
		System.out.println("Put in " + value);
		value = normalized(value, offset);
		if (value == -1) {
			String errInfo = String.format("非法的值[%d], 值必须在[%d, %d]之间", max, offset+1, offset+LimitMaxValue);
			throw new IllegalArgumentException(errInfo);
		}
		if (value > max) {
			grow(value);
		}
		int index = (int)(value >>> 3);
		byte mod = (byte)(value & 0x07);
		if (mod == 0) {
			index -= 1;
		}
		byte mask = mask(mod);
		
		buffer.position(index);
		byte origin = buffer.get();
		buffer.position(index);
		buffer.put((byte)(origin | mask));
	}
	
	public void remove(long value) {
		value = normalized(value, offset);
		if (value == -1) {
			return;
		}
		int index = (int)(value >>> 3);
		byte mod = (byte)(value & 0x07);
		if (mod == 0) {
			index -= 1;
		}
		byte mask = mask(mod);

		buffer.position(index);
		byte origin = buffer.get();
		buffer.position(index);
		buffer.put((byte)(origin & ~mask));
	}
	
	public static void main(String[] args) {
		BitMap bigMap = new BitMap(10, 0); //range from (offset + 1) to (offset + 2^32) 
		bigMap.put(6);
		bigMap.put(24);
		bigMap.put(200000000);
		System.out.println(bigMap.existed(6));
		System.out.println(bigMap.existed(23));
		System.out.println(bigMap.existed(24));
		bigMap.remove(24);
		System.out.println(bigMap.existed(24));
		System.out.println(bigMap.existed(200000000));
	}
}
