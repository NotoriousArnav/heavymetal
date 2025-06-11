import { clsx, type ClassValue } from "clsx"
import { twMerge } from "tailwind-merge"
import axios from "axios"

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs))
}

// Make a Class to interact with the API
//export class HeavymetalAPI {
//  /*
//   * Make requests to the Heavymetal API whose base URL is taken from the environment variable
//   * VITE_HEAVYMETAL_API_URL.
//   * */
//  private static baseURL = import.meta.env.VITE_HEAVYMETAL_API_URL || "http://127.0.0.1:8080"
//  private static axiosInstance = axios.create({
//    baseURL: HeavymetalAPI.baseURL,
//    headers: {
//      "Content-Type": "application/json",
//    },
//  })
//
//  public static async get<T>(path: string, params?: Record<string, any>): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.get<T>(path, { params })
//      return response.data
//    } catch (error) {
//      console.error("Error in GET request:", error)
//      throw error
//    }
//  }
//
//  public static async post<T>(path: string, data?: Record<string, any>): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.post<T>(path, data)
//      return response.data
//    } catch (error) {
//      console.error("Error in POST request:", error)
//      throw error
//    }
//  }
//
//  public static async put<T>(path: string, data?: Record<string, any>): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.put<T>(path, data)
//      return response.data
//    } catch (error) {
//      console.error("Error in PUT request:", error)
//      throw error
//    }
//  }
//
//  public static async delete<T>(path: string): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.delete<T>(path)
//      return response.data
//    } catch (error) {
//      console.error("Error in DELETE request:", error)
//      throw error
//    }
//  }
//
//  public static async patch<T>(path: string, data?: Record<string, any>): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.patch<T>(path, data)
//      return response.data
//    } catch (error) {
//      console.error("Error in PATCH request:", error)
//      throw error
//    }
//  }
//
//  public static async head<T>(path: string, params?: Record<string, any>): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.head<T>(path, { params })
//      return response.data
//    } catch (error) {
//      console.error("Error in HEAD request:", error)
//      throw error
//    }
//  }
//
//  public static async options<T>(path: string, params?: Record<string, any>): Promise<T> {
//    try {
//      const response = await HeavymetalAPI.axiosInstance.options<T>(path, { params })
//      return response.data
//    } catch (error) {
//      console.error("Error in OPTIONS request:", error)
//      throw error
//    }
//  }
//  
//  public async getTrackFromID<T>(trackId: string): Promise<T> {
//    return HeavymetalAPI.get<T>(`/track/${trackId}`)
//  }
//
//  public async searchTracksFuzzy<T>(query: string): Promise<T> {
//    return HeavymetalAPI.get<T>(`/search/track/${query}`)
//  }
//
//}
