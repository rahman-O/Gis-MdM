import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import axios from 'axios'
import { useAuth } from '@/shared/hooks/useAuth'

const loginSchema = z.object({
  login: z.string().min(1, 'Username is required'),
  password: z.string().min(1, 'Password is required'),
})

export type LoginFormValues = z.infer<typeof loginSchema>

export function useLogin() {
  const { login } = useAuth()
  const [isLoading, setIsLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const form = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { login: '', password: '' },
  })

  const onSubmit = async (values: LoginFormValues) => {
    setIsLoading(true)
    setError(null)
    try {
      await login(values)
    } catch (err: unknown) {
      if (axios.isAxiosError(err)) {
        if (!err.response) {
          setError('Cannot reach the server. Make sure the backend is running on port 8080.')
        } else if (err.response.status === 401 || err.response.status === 403) {
          setError('Invalid username or password.')
        } else {
          setError(`Server error (${err.response.status}). Please try again.`)
        }
      } else {
        setError(err instanceof Error ? err.message : 'An unexpected error occurred.')
      }
    } finally {
      setIsLoading(false)
    }
  }

  return { form, isLoading, error, onSubmit: form.handleSubmit(onSubmit) }
}
