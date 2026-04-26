declare module 'jsencrypt' {
  export default class JSEncrypt {
    constructor();
    setPublicKey(key: string): void;
    encrypt(plaintext: string): string | false;
    setPrivateKey(key: string): void;
    decrypt(ciphertext: string): string | false;
  }
}
