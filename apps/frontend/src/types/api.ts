export interface User {
  user_id: number;
  firstname: string;
  lastname: string;
  email: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RegisterRequest {
  firstname: string;
  lastname: string;
  email: string;
  password: string;
  confirm_password: string;
}

export interface Card {
  card_id: number;
  user_id: number;
  card_number: string;
  card_type: string;
  expire_date: string;
  cvv: string;
  card_provider: string;
}

export interface Saldo {
  saldo_id: number;
  card_number: string;
  total_balance: number;
}

export interface TopupRequest {
  card_number: string;
  topup_no: string;
  topup_amount: number;
  topup_method: string;
}

export interface TransferRequest {
  transfer_from: string;
  transfer_to: string;
  transfer_amount: number;
}

export interface WithdrawRequest {
  card_number: string;
  withdraw_amount: number;
  withdraw_time: string;
}

export interface Transaction {
  transaction_id: number;
  card_number: string;
  amount: number;
  payment_method: string;
  merchant_id: number;
  transaction_time: string;
}
