import { fromBech32 } from "@cosmjs/encoding";
import { DirectSecp256k1HdWallet } from "@cosmjs/proto-signing";
import CryptoJS from "crypto-js";
import { defineStore } from "pinia";

import { useClient } from "@/composables/useClient";
import type { EncodedWallet, Nullable, Wallet } from "@/utils/interfaces";

const DEV_MNEMONIC_STORAGE_KEY = "ferac.dev.mnemonic";
const DEV_WALLET_NAME_STORAGE_KEY = "ferac.dev.walletName";
const DEV_WALLETS_STORAGE_KEY = "ferac.dev.wallets";

const getMnemonicStorageKey = () =>
  `ferac-dev-wallet:${window.location.origin}`;

const sameAddress = (firstAddress: string, secondAddress: string) => {
  try {
    const first = fromBech32(firstAddress).data;
    const second = fromBech32(secondAddress).data;
    return (
      first.length === second.length &&
      first.every((byte, index) => byte === second[index])
    );
  } catch {
    return firstAddress === secondAddress;
  }
};

type StoredDevWallet = {
  name: string;
  address: string;
  encryptedMnemonic: string;
};

const readStoredDevWallets = (): StoredDevWallet[] => {
  const stored = window.localStorage.getItem(DEV_WALLETS_STORAGE_KEY);
  if (!stored) {
    return [];
  }

  try {
    return JSON.parse(stored) as StoredDevWallet[];
  } catch {
    return [];
  }
};

const writeStoredDevWallets = (wallets: StoredDevWallet[]) => {
  window.localStorage.setItem(DEV_WALLETS_STORAGE_KEY, JSON.stringify(wallets));
};

export const useWalletStore = defineStore("wallet", {
  state: () => ({
    wallets:
      (JSON.parse(
        window.localStorage.getItem("wallets") ?? "null"
      ) as Array<EncodedWallet>) || ([] as Array<EncodedWallet>),
    activeWallet: null as Nullable<Wallet>,
    activeClient: null as Nullable<ReturnType<typeof useClient>>,
    devWallets: readStoredDevWallets(),
    selectedAddress: "",
    authorized: false,
    backupState: true,
    gasPrice: "0.0000025token",
  }),
  getters: {
    getClient: (state) => state.activeClient,
    getGasPrice: (state) => state.gasPrice,
    getWallet: (state) => state.activeWallet,
    getDevWallets: (state) => state.devWallets,
    getAddress: (state) => state.selectedAddress,
    getPath: (state) => {
      if (state.activeWallet && state.activeWallet.HDpath) {
        return (
          state.activeWallet.HDpath +
          state.activeWallet.accounts.find(
            (x) => x.address == state.selectedAddress
          )?.pathIncrement
        );
      } else {
        return null;
      }
    },
    getNameAvailable: (state) => (name: string) => {
      return state.wallets.findIndex((x) => x.name == name) == -1;
    },
    getLastWallet: (state) => {
      if (state.activeWallet) {
        return state.activeWallet.name;
      } else {
        return window.localStorage.getItem("lastWallet");
      }
    },
    getLoggedIn: (state) => state.activeClient !== null,
    getSigner: (state) => {
      if (state.activeClient) {
        return state.activeClient.signer;
      } else {
        return null;
      }
    },
  },
  actions: {
    async signOut() {
      const activeAddress =
        this.selectedAddress || this.activeWallet?.accounts[0]?.address || "";
      if (activeAddress) {
        this.devWallets = this.devWallets.filter(
          (wallet) => !sameAddress(wallet.address, activeAddress)
        );
      }
      writeStoredDevWallets(this.devWallets);
      const nextWallet = this.devWallets[0];
      window.localStorage.removeItem("lastWallet");
      window.localStorage.removeItem(DEV_MNEMONIC_STORAGE_KEY);
      window.localStorage.removeItem(DEV_WALLET_NAME_STORAGE_KEY);
      this.selectedAddress = "";
      this.activeClient?.removeSigner();
      this.activeClient = null;
      this.activeWallet = null;
      this.authorized = false;
      if (nextWallet) {
        await this.switchToMnemonicWallet(nextWallet.address);
      }
    },
    async connectWithKeplr() {
      const client = useClient();

      try {
        const wallet: Wallet = {
          name: "Keplr Integration",
          mnemonic: null,
          HDpath: null,
          password: null,
          prefix: client.env.prefix ?? "cosmos",
          pathIncrement: null,
          accounts: [],
        };
        await client.useKeplr();
        if (!client.signer) {
          throw new Error("Keplr signer not available");
        }
        const [account] = await client.signer.getAccounts();
        wallet.accounts.push({ address: account.address, pathIncrement: null });

        this.activeWallet = wallet;
        this.selectedAddress = account.address;
        window.localStorage.setItem("lastWallet", wallet.name);
        if (
          this.activeWallet &&
          this.activeWallet.name &&
          this.activeWallet.password
        ) {
          this.wallets.push({
            name: this.activeWallet.name,
            wallet: CryptoJS.AES.encrypt(
              JSON.stringify(this.activeWallet),
              this.activeWallet.password
            ).toString(),
          });
        }
        if (
          this.activeWallet.name == "Keplr Integration" &&
          !this.activeWallet.password
        ) {
          this.wallets.push({
            name: this.activeWallet.name,
            wallet: JSON.stringify(this.activeWallet),
          });
        }

        this.activeClient = client;
      } catch (e) {
        console.error(e);
      }
      this.storeWallets();
    },
    async connectWithMnemonic(mnemonic: string, name = "Local Dev Wallet") {
      const client = useClient();
      const walletSigner = await DirectSecp256k1HdWallet.fromMnemonic(
        mnemonic.trim(),
        { prefix: client.env.prefix ?? "cosmos" }
      );
      const [account] = await walletSigner.getAccounts();
      const wallet: Wallet = {
        name,
        mnemonic: null,
        HDpath: null,
        password: null,
        prefix: client.env.prefix ?? "cosmos",
        pathIncrement: null,
        accounts: [{ address: account.address, pathIncrement: null }],
      };

      client.useSigner(walletSigner);
      this.selectedAddress = account.address;
      this.activeWallet = wallet;
      this.activeClient = client;
      this.authorized = true;
      const encryptedMnemonic = CryptoJS.AES.encrypt(
        mnemonic.trim(),
        getMnemonicStorageKey()
      ).toString();
      const storedWallets = this.devWallets.filter(
        (storedWallet) => storedWallet.address !== account.address
      );
      storedWallets.push({
        name,
        address: account.address,
        encryptedMnemonic,
      });
      this.devWallets = storedWallets;
      writeStoredDevWallets(storedWallets);
      window.localStorage.setItem(DEV_MNEMONIC_STORAGE_KEY, encryptedMnemonic);
      window.localStorage.setItem(DEV_WALLET_NAME_STORAGE_KEY, name);
    },
    async restoreMnemonicWallet() {
      const storedWallets = this.devWallets;
      let selectedWallet = storedWallets.find(
        (storedWallet) =>
          storedWallet.name === window.localStorage.getItem("lastWallet")
      );

      const legacyMnemonic = window.localStorage.getItem(
        DEV_MNEMONIC_STORAGE_KEY
      );
      if (!selectedWallet && legacyMnemonic && storedWallets.length === 0) {
        const mnemonic = CryptoJS.AES.decrypt(
          legacyMnemonic,
          getMnemonicStorageKey()
        ).toString(CryptoJS.enc.Utf8);
        if (mnemonic) {
          await this.connectWithMnemonic(
            mnemonic,
            window.localStorage.getItem(DEV_WALLET_NAME_STORAGE_KEY) ||
              "Local Dev Wallet"
          );
          selectedWallet = this.devWallets[0];
        }
      }

      selectedWallet ||= storedWallets[0];
      if (!selectedWallet) {
        return false;
      }

      const encryptedMnemonic = selectedWallet.encryptedMnemonic;
      if (!encryptedMnemonic) {
        return false;
      }

      const mnemonic = CryptoJS.AES.decrypt(
        encryptedMnemonic,
        getMnemonicStorageKey()
      ).toString(CryptoJS.enc.Utf8);
      if (!mnemonic) {
        throw new Error("Saved wallet data is invalid");
      }

      await this.connectWithMnemonic(mnemonic, selectedWallet.name);
      return true;
    },
    async switchToMnemonicWallet(address: string) {
      const storedWallet = this.devWallets.find(
        (wallet) => wallet.address === address
      );
      if (!storedWallet) {
        throw new Error("Wallet not found");
      }

      const mnemonic = CryptoJS.AES.decrypt(
        storedWallet.encryptedMnemonic,
        getMnemonicStorageKey()
      ).toString(CryptoJS.enc.Utf8);
      if (!mnemonic) {
        throw new Error("Saved wallet data is invalid");
      }

      await this.connectWithMnemonic(mnemonic, storedWallet.name);
    },
    storeWallets() {
      window.localStorage.setItem("wallets", JSON.stringify(this.wallets));
      this.backupState = false;
    },
  },
});
