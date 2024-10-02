import { useState } from "react";
import { JSEncrypt } from "jsencrypt";

export default function useRSA(publicKey: string) {
  const [encryptedData, setEncryptedData] = useState('');
  const [rsaEncError, setHasError] = useState(false);

  const encrypt = (data: string) => {
    try {
      const encrypt = new JSEncrypt();
      encrypt.setPublicKey(publicKey);
      const encrypted = encrypt.encrypt(data);
      
      if (!encrypted) {
        setEncryptedData('');
        setHasError(true);
        return;
      }
      
      setEncryptedData(encrypted);
      setHasError(false);
      return encrypted;
    } catch (err) {
      setEncryptedData('');
      setHasError(true);
    }
  };

  // ToDo: Add decryption function

  return { encryptedData, rsaEncError, encrypt };
};
