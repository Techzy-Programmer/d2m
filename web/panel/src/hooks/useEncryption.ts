import { useState } from "react";
import { JSEncrypt } from "jsencrypt";
import CryptoJS from 'crypto-js';

interface EncryptionResult {
  encData: string;
  encKey: string;
}

export default function useEncryption(publicKey: string) {
  const [encryptedResult, setEncryptedResult] = useState<EncryptionResult | null>(null);
  const [encError, setEncError] = useState(false);

  // Function to generate a random AES key (32 bytes = 256 bits)
  const generateAESKey = (length: number = 32): string => {
    const charset = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()_+-=[]{}|;:,.<>?';
    const randomValues = new Uint8Array(length);
    window.crypto.getRandomValues(randomValues);  // Secure random values
    
    let result = '';
    for (let i = 0; i < length; i++) {
      result += charset.charAt(randomValues[i] % charset.length);
    }
    
    return result;
  };

  // Function to encrypt data with AES-256-CBC
  const encryptWithAES = (data: string, key: string): string => {
    const iv = CryptoJS.lib.WordArray.random(16);
    
    const encrypted = CryptoJS.AES.encrypt(data, CryptoJS.enc.Utf8.parse(key), {
      iv: iv,
      mode: CryptoJS.mode.CBC,
      padding: CryptoJS.pad.Pkcs7
    });
    
    const ivAndData = iv.concat(encrypted.ciphertext);
    return ivAndData.toString(CryptoJS.enc.Base64);
  };

  // Function to encrypt the AES key using RSA with the provided public key
  const encryptAESKeyWithRSA = (aesKey: string): string | null => {
    try {
      const rsaEncrypt = new JSEncrypt();
      rsaEncrypt.setPublicKey(publicKey);
      const encryptedKey = rsaEncrypt.encrypt(aesKey);
      
      if (!encryptedKey) {
        throw new Error('RSA encryption failed');
      }
      return encryptedKey;
    } catch (error) {
      console.error('RSA encryption error:', error);
      return null;
    }
  };

  // Main encryption function
  const encrypt = (data: string): EncryptionResult | null => {
    try {
      const aesKey = generateAESKey();
      const encData = encryptWithAES(data, aesKey);
      const encKey = encryptAESKeyWithRSA(aesKey);
      
      if (!encKey) {
        setEncError(true);
        setEncryptedResult(null);
        return null;
      }
      
      const result = {
        encData,
        encKey
      };
      
      setEncryptedResult(result);
      setEncError(false);
      return result;
    } catch (error) {
      console.error('Encryption error:', error);
      setEncryptedResult(null);
      setEncError(true);
      return null;
    }
  };

  return {
    encryptedResult,
    encError,
    encrypt,
  };
}
