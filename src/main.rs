use rust_fzf;
use std::{fs, io};

fn get_projects(root_dir: String) -> io::Result<()> {
    let paths = fs::read_dir(root_dir)?.filter(|path| path.is_dir())
                .filter(|path| !path.starts_with(".")).collect::<Vec<String>>();
    Ok()
}
fn main() {
    let paths = fs::read_dir("./")
        .unwrap()
        .map(|path| path.unwrap().path())
        .map(|path| fs::canonicalize(path).unwrap())
        .map(|path| path.to_str().unwrap().to_string())
        .collect::<Vec<String>>();

    let project = rust_fzf::select(paths, vec![]);
    println!("{}", project);
}
