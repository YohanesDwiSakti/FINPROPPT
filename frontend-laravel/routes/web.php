<?php

use App\Models\User;
use Illuminate\Http\Request;
use Illuminate\Support\Facades\DB;
use Illuminate\Support\Facades\Hash;
use Illuminate\Support\Facades\Route;
use Illuminate\Validation\ValidationException;

Route::get('/', function () {
    if (!session('role')) {
        return redirect('/login')->with('message', 'Login dulu untuk membuka aplikasi.');
    }

    if (session('role') === 'admin') {
        return redirect('/admin');
    }

    return view('welcome');
});

Route::get('/admin', function () {
    if (session('role') !== 'admin') {
        return redirect('/login')->with('message', 'Silakan login sebagai admin untuk membuka dashboard.');
    }

    return view('admin');
});

Route::get('/login', function () {
    if (session('role') === 'admin') {
        return redirect('/admin');
    }

    if (session('role') === 'customer') {
        return redirect('/');
    }

    return view('auth', ['mode' => 'login']);
});

Route::post('/login', function (Request $request) {
    $request->validate([
        'email' => ['required', 'email'],
        'password' => ['required'],
        'role' => ['required', 'in:customer,admin'],
    ]);

    $user = User::where('email', $request->input('email'))
        ->where('role', $request->input('role'))
        ->first();

    if (! $user || ! Hash::check($request->input('password'), $user->password)) {
        throw ValidationException::withMessages([
            'email' => 'Email, password, atau role tidak sesuai.',
        ]);
    }

    session([
        'user_id' => $user->id,
        'user_name' => $user->name,
        'role' => $user->role,
    ]);

    return $user->role === 'admin'
        ? redirect('/admin')
        : redirect('/');
});

Route::get('/register', function () {
    if (session('role') === 'admin') {
        return redirect('/admin');
    }

    if (session('role') === 'customer') {
        return redirect('/');
    }

    return view('auth', ['mode' => 'register']);
});

Route::post('/register', function (Request $request) {
    $request->validate([
        'name' => ['required'],
        'email' => ['required', 'email', 'unique:users,email'],
        'password' => ['required', 'min:6'],
    ]);

    $user = User::create([
        'name' => $request->input('name'),
        'email' => $request->input('email'),
        'password' => Hash::make($request->input('password')),
        'role' => 'customer',
    ]);

    session([
        'user_id' => $user->id,
        'user_name' => $user->name,
        'role' => 'customer',
    ]);

    return redirect('/');
});

Route::post('/logout', function () {
    session()->flush();

    return redirect('/');
});

Route::post('/manifests', function (Request $request) {
    if (session('role') !== 'admin') {
        return response()->json(['message' => 'Unauthorized'], 403);
    }

    $data = $request->validate([
        'receipt' => ['required'],
        'status' => ['required'],
        'location' => ['nullable'],
    ]);

    DB::table('manifests')->updateOrInsert(
        ['receipt' => strtoupper(trim($data['receipt']))],
        [
            'status' => $data['status'],
            'location' => $data['location'] ?? null,
            'updated_by' => session('user_id'),
            'updated_at' => now(),
            'created_at' => now(),
        ],
    );

    return response()->json([
        'message' => 'Manifest ' . strtoupper(trim($data['receipt'])) . ' berhasil disimpan ke database.',
    ]);
});
